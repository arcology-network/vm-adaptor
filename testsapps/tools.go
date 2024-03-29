package tests

import (
	"errors"
	"math"
	"math/big"

	"github.com/arcology-network/common-lib/cachedstorage"
	commontypes "github.com/arcology-network/common-lib/types"
	concurrenturl "github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/interfaces"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	"github.com/arcology-network/evm/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	evmcore "github.com/arcology-network/evm/core"
	evmcoretypes "github.com/arcology-network/evm/core/types"

	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"

	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/execution"
)

func MainTestConfig() *execution.Config {
	vmConfig := vm.Config{}
	cfg := &execution.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: big.NewInt(0),
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(0),
		Coinbase:    &eucommon.Coinbase,
		GasLimit:    math.MaxUint64, // Should come from the message
		Difficulty:  big.NewInt(0),
	}
	cfg.Chain = new(execution.DummyChain)
	return cfg
}

func NewTestEU() (*execution.EU, *execution.Config, interfaces.Datastore, *concurrenturl.ConcurrentUrl, []interfaces.Univalue) {
	persistentDB := cachedstorage.NewDataStore(nil, cachedstorage.NewCachePolicy(0, 1), cachedstorage.NewMemDB(), ccurlstorage.Rlp{}.Encode, ccurlstorage.Rlp{}.Decode)
	persistentDB.Inject(ccurlcommon.ETH10_ACCOUNT_PREFIX, commutative.NewPath())
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccapi.NewAPI(url)

	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)

	statedb.CreateAccount(eucommon.Alice)
	statedb.AddBalance(eucommon.Alice, new(big.Int).SetUint64(1e18))

	statedb.CreateAccount(eucommon.Bob)
	statedb.AddBalance(eucommon.Bob, new(big.Int).SetUint64(1e18))

	// statedb.CreateAccount(eucommon.RUNTIME_HANDLER)
	// statedb.AddBalance(eucommon.RUNTIME_HANDLER, new(big.Int).SetUint64(1e18))

	_, transitions := api.StateFilter().ByType()
	// indexer.Univalues(transitionsFiltered).Print()

	// fmt.Println("\n" + eucommon.FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{0})
	api = ccapi.NewAPI(url)
	statedb = eth.NewImplStateDB(api)

	config := MainTestConfig()
	config.Coinbase = &eucommon.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	return execution.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config, db, url, transitions
}

func DepolyContract(eu *execution.EU, config *execution.Config, code string, funcName string, inputData []byte, nonce uint64, checkNonce bool) (error, *execution.Config, *execution.EU, *evmcoretypes.Receipt) {
	msg := core.NewMessage(eucommon.Alice, nil, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false)
	stdMsg := &execution.StandardMessage{
		ID:     1,
		TxHash: [32]byte{1, 1, 1},
		Native: &msg, // Build the message
		Source: commontypes.TX_SOURCE_LOCAL,
	}

	receipt, result, err := eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it

	if result.Err != nil {
		return result.Err, config, eu, nil
	}

	if err != nil || receipt.Status != 1 {
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		return errors.New("Error: Deployment failed!!!" + errmsg), config, eu, nil
	}

	_, transitionsFiltered := eu.Api().StateFilter().ByType()
	eu.Api().Ccurl().Import(transitionsFiltered)
	eu.Api().Ccurl().Sort()
	eu.Api().Ccurl().Commit([]uint32{1})

	return nil, config, eu, receipt
}

func DeployThenInvoke(targetPath, file, version, contractName, funcName string, inputData []byte, checkNonce bool) (error, *execution.EU, *evmcoretypes.Receipt) {
	code, err := compiler.CompileContracts(targetPath, file, version, contractName, false)
	eu, config, _, _, _ := NewTestEU()
	if err != nil || len(code) == 0 {
		return err, nil, nil
	}

	err, _, eu, receipt := DepolyContract(eu, config, code, funcName, inputData, 0, checkNonce)

	if len(funcName) == 0 || err != nil {
		return err, eu, receipt
	}

	data := crypto.Keccak256([]byte(funcName))[:4]
	data = append(data, inputData...)
	err, eu, execResult, receipt := CallContract(eu, receipt.ContractAddress, data, 0, checkNonce)

	if err != nil || receipt.Status != 1 {
		return execResult.Err, eu, receipt
	}

	if execResult != nil && execResult.Err != nil {
		return execResult.Err, eu, receipt
	}
	return nil, eu, receipt
}

func CallContract(eu *execution.EU, contractAddress common.Address, inputData []byte, nonceIncrement uint64, checkNonce bool) (error, *execution.EU, *evmcore.ExecutionResult, *evmcoretypes.Receipt) {
	// data := crypto.Keccak256([]byte(funcName))[:4]
	// inputData = append(data, inputData...)

	msg := core.NewMessage(eucommon.Alice, &contractAddress, 10+nonceIncrement, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), inputData, nil, false)
	stdMsg := &execution.StandardMessage{
		ID:     1,
		TxHash: [32]byte{1, 1, 1},
		Native: &msg, // Build the message
		Source: commontypes.TX_SOURCE_LOCAL,
	}

	config := MainTestConfig()
	config.Coinbase = &eucommon.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	var execResult *evmcore.ExecutionResult
	receipt, execResult, err := eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it
	// _, transitions := eu.Api().StateFilter().ByType()

	// msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	// receipt, execResult, _ := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
	// _, transitions = eu.Api().StateFilter().ByType()

	if err != nil {
		return nil, nil, execResult, receipt
	}

	if receipt.Status != 1 {
		return execResult.Err, eu, execResult, receipt
	}

	if execResult != nil && execResult.Err != nil {
		return execResult.Err, eu, execResult, receipt
	}
	return nil, eu, execResult, receipt
}
