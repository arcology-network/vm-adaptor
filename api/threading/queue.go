package multiprocess

import (
	"crypto/sha256"
	"math/big"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	"github.com/arcology-network/concurrenturl/univalue"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	types "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"

	cceu "github.com/arcology-network/vm-adaptor"

	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

type Job struct {
	sender      evmcommon.Address
	caller      evmcommon.Address
	callee      evmcommon.Address
	message     types.Message
	receipt     *types.Receipt
	result      *core.ExecutionResult
	prechkErr   error
	accesses    []ccurlcommon.UnivalueInterface
	transitions []ccurlcommon.UnivalueInterface
	ccurl       *concurrenturl.ConcurrentUrl
}

// APIs under the concurrency namespace
type Queue struct {
	jobs       []*Job
	arbitrator *indexer.ArbitratorSlow
}

func NewJobQueue() *Queue {
	return &Queue{
		jobs:       []*Job{},
		arbitrator: indexer.NewArbitratorSlow(),
	}
}

func (this *Queue) Length() uint64 { return uint64(len(this.jobs)) }

func (this *Queue) At(idx uint64) *Job {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *Job { return this.jobs[idx] }, nil)
}

func (this *Queue) Del(idx uint64) {
	common.IfThenDo(idx < uint64(len(this.jobs)), func() { this.jobs[idx] = nil }, func() {})
	common.RemoveIf(&this.jobs, func(job *Job) bool { return job == nil })
}

func (this *Queue) Add(origin, calleeAddr evmcommon.Address, funCall []byte) int {
	this.jobs = append(this.jobs,
		&Job{
			sender: origin,
			caller: evmcommon.Address{},
			callee: calleeAddr,
			message: types.NewMessage( // Build the message
				origin,
				&calleeAddr,
				0,
				new(big.Int).SetUint64(0), // Amount to transfer
				1e15,
				new(big.Int).SetUint64(1),
				funCall,
				nil,
				false, // Don't checking nonce
			),
		},
	)
	return len(this.jobs) - 1
}

func (this *Queue) Snapshot(mainApiRouter eucommon.ConcurrentApiRouterInterface) ccurlcommon.DatastoreInterface {
	transitions := mainApiRouter.Ccurl().Export()                                                         // Get the all up-to-date transitions from the main thread
	mainProcessTrans := univalue.Univalues(common.Clone(transitions)).To(univalue.TransitionFilters()...) // Filter out unwanted ones

	transientDB := ccurlstorage.NewTransientDB(mainApiRouter.Ccurl().WriteCache().Store()) // Should be the same as Importer().Store()
	snapshot := concurrenturl.NewConcurrentUrl(transientDB).Import(mainProcessTrans).Sort()
	return snapshot.Commit([]uint32{mainApiRouter.TxIndex()}).Importer().Store() // Commit these changes to the a transient DB
}

func (this *Queue) Run(threads uint8, mainApiRouter eucommon.ConcurrentApiRouterInterface) bool {
	snapshot := this.Snapshot(mainApiRouter)
	// _, _1 := this.apiRouter.Ccurl().ExportAll()
	// univalue.Univalues(_1).Print()

	// t0 := time.Now()
	config := cceu.NewConfig().SetCoinbase(mainApiRouter.Coinbase()) // Share the same coinbase as the main thread
	// executor := func(start, end, index int, args ...interface{}) {
	// for i := start; i < end; i++ {
	for i := 0; i < len(this.jobs); i++ {
		ccurl := (&concurrenturl.ConcurrentUrl{}).New(
			indexer.NewWriteCache(snapshot, mainApiRouter.Ccurl().Platform),
			mainApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

		txHash := sha256.Sum256(codec.Uint64(i).Encode())

		this.jobs[i].ccurl = ccurl
		statedb := eth.NewImplStateDB(ccurl)   // Eth state DB
		statedb.Prepare(txHash, [32]byte{}, i) // tx hash , block hash and tx index

		apiRounter := mainApiRouter.New(txHash, uint32(i), ccurl)
		eu := cceu.NewEU(
			params.MainnetChainConfig,
			vm.Config{},
			statedb,
			apiRounter, // Tx hash, tx id and url
		)

		this.jobs[i].accesses, this.jobs[i].transitions, this.jobs[i].receipt, this.jobs[i].result, this.jobs[i].prechkErr =
			eu.Run(eu.Api().TxHash(), int(eu.Api().TxIndex()), &this.jobs[i].message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.jobs[i].message))
	}
	// }
	// common.ParallelWorker(len(this.jobs), int(threads), executor)
	// fmt.Println("Run: ", time.Since(t0))

	// Put all the access records together
	length := 0
	common.Foreach(this.jobs, func(job **Job) { length += len((*(*job)).accesses) }) // Pre-allocation

	accesseVec := make([]ccurlcommon.UnivalueInterface, 0, length)
	common.Foreach(this.jobs, func(job **Job) { accesseVec = append(accesseVec, (*(*job)).accesses...) })

	// Detect potential conflicts}
	_, conflicTxs := indexer.NewArbitratorSlow().Detect(accesseVec)
	univalue.Univalues(accesseVec).Print()

	// Clear up conflicting txs and their state changes
	jobs := common.SetIndices(common.Clone(this.jobs), conflicTxs, func(job *Job) *Job { return nil })
	common.RemoveIf(&jobs, func(job *Job) bool { return job == nil }) // Should not remove the jobs directly

	// t0 = time.Now()
	this.WriteBack(mainApiRouter, jobs) // Merge back to the main write cache
	// fmt.Println("Commit: ", time.Since(t0))
	return true
}

// Merge all the transitions back to the main cache
func (this *Queue) WriteBack(mainApiRouter eucommon.ConcurrentApiRouterInterface, jobs []*Job) {
	for i := 0; i < len(jobs); i++ {
		mainApiRouter.Ccurl().WriteCache().MergeFrom(jobs[i].ccurl.WriteCache(), mainApiRouter.TxIndex())
	}
}

func (this *Queue) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
