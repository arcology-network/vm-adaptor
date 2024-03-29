package api

import (
	"math"
	"sync/atomic"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"
	"github.com/holiman/uint256"

	"github.com/arcology-network/vm-adaptor/abi"
	execution "github.com/arcology-network/vm-adaptor/execution"

	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type MultiprocessHandlers struct {
	*BaseHandlers
	erros   []error
	jobseqs []*execution.JobSequence
}

func NewMultiprocessHandlers(ethApiRouter eucommon.EthApiRouter) *MultiprocessHandlers {
	handler := &MultiprocessHandlers{
		erros:   []error{},
		jobseqs: []*execution.JobSequence{},
	}
	handler.BaseHandlers = NewBaseHandlers(ethApiRouter, handler)
	return handler
}

func (this *MultiprocessHandlers) Address() [20]byte { return eucommon.MULTIPROCESS_HANDLER }

func (this *MultiprocessHandlers) Run(caller [20]byte, input []byte) ([]byte, bool, int64) {
	if atomic.AddUint64(&eucommon.TotalSubProcesses, 1); !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, 0
	}

	input, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt64)
	if err != nil {
		return []byte{}, false, 0
	}

	numThreads, err := abi.DecodeTo(input, 0, uint64(1), 1, 8)
	if err != nil {
		return []byte{}, false, 0
	}
	threads := common.Min(common.Max(uint8(numThreads), 1), math.MaxUint8) // [1, 255]

	path := this.Connector().Key(caller)
	length, successful, fee := this.Length(path)
	length = common.Min(eucommon.MAX_VM_INSTANCES, length)

	if !successful {
		return []byte{}, successful, fee
	}

	generation := execution.NewGeneration(0, threads, []*execution.JobSequence{})
	fees := make([]int64, length)
	this.erros = make([]error, length)
	this.jobseqs = make([]*execution.JobSequence, length)

	for i := uint64(0); i < length; i++ {
		funCall, successful, fee := this.GetByIndex(path, uint64(i))
		if fees[i] = fee; successful {
			this.jobseqs[i], this.erros[i] = this.toJobSeq(funCall)
		}
		generation.Add(this.jobseqs[i])
	}
	transitions := generation.Run(this.Api())

	// Sub processes may have been spawned during the execution, recheck it.
	if !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, fee
	}

	// Unify tx IDs c
	mainTxID := uint32(this.Api().GetEU().(*execution.EU).Message().ID)
	common.Foreach(transitions, func(v *interfaces.Univalue, _ int) { (*v).SetTx(mainTxID) })

	this.Api().Ccurl().WriteCache().AddTransitions(transitions) // Merge the write cache to the main cache
	return []byte{}, true, common.Sum[int64](fees)
}

func (this *MultiprocessHandlers) toJobSeq(input []byte) (*execution.JobSequence, error) {
	gasLimit, value, calleeAddr, funCall, err := abi.Parse4(input,
		uint64(0), 1, 32,
		uint256.NewInt(0), 1, 32,
		[20]byte{}, 1, 32,
		[]byte{}, 2, math.MaxInt64)

	if err != nil {
		return nil, err
	}

	newJobSeq := &execution.JobSequence{
		ID:        uint32(this.BaseHandlers.Api().GetSerialNum(eucommon.SUB_PROCESS)),
		ApiRouter: this.BaseHandlers.Api(),
	}

	transfer := value.ToBig()
	addr := evmcommon.Address(calleeAddr)
	evmMsg := evmcore.NewMessage( // Build the message
		this.BaseHandlers.Api().Origin(),
		&addr,
		0,
		transfer, // Amount to transfer
		gasLimit,
		this.BaseHandlers.Api().GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	stdMsg := &execution.StandardMessage{
		ID:     uint64(newJobSeq.ID), // this is the problem !!!!
		Native: &evmMsg,
		TxHash: newJobSeq.DeriveNewHash(this.BaseHandlers.Api().GetEU().(*execution.EU).Message().TxHash),
	}

	newJobSeq.StdMsgs = []*execution.StandardMessage{stdMsg}
	return newJobSeq, nil
}
