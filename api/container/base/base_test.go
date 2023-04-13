package base

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	ccurlcommon "github.com/arcology-network/concurrenturl/v2/common"
	ccurlstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	cceu "github.com/arcology-network/vm-adaptor"
	ccapi "github.com/arcology-network/vm-adaptor/api"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

func TestBase(t *testing.T) {
	config := compiler.MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(ccurlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(ccurlcommon.NewPlatform().Eth10Account(), meta)
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(compiler.Coinbase)
	statedb.CreateAccount(compiler.User1)
	statedb.AddBalance(compiler.User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + compiler.FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api := ccapi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := cceu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &compiler.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	pyCompiler := filepath.Dir(filepath.Dir(filepath.Dir(currentPath))) + "/compiler/compiler.py"
	code, err := compiler.CompileContracts(pyCompiler, "./base_test.sol", "BaseTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(compiler.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)     // Build the message
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + compiler.FormatTransitions(transitions))
	// t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== Call length() ==================================
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	errs := url.Commit([]uint32{1})
	if len(errs) != 0 {
		t.Error(errs)
		return
	}
	api = ccapi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu = cceu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.BlockNumber = new(big.Int).SetUint64(10000001)
	config.Time = new(big.Int).SetUint64(10000001)

	data := crypto.Keccak256([]byte("length()"))[:4]
	data = append(data, evmcommon.BytesToHash(compiler.User1.Bytes()).Bytes()...)
	data = append(data, evmcommon.BytesToHash([]byte{0xcc}).Bytes()...)
	msg = types.NewMessage(compiler.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	_, transitions, receipt, err = eu.Run(evmcommon.BytesToHash([]byte{2, 2, 2}), 2, &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg))
	t.Log("\n" + compiler.FormatTransitions(transitions))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Failed to calll length()!!!", err)
	}

	// Get.
	// url = concurrenturl.NewConcurrentUrl(db)
	// url.Import(transitions)
	// url.PostImport()
	// errs = url.Commit([]uint32{2})
	// if len(errs) != 0 {
	// 	t.Error(errs)
	// 	return
	// }
	// api = ccapi.NewAPI(url)
	// statedb = eth.NewImplStateDB(url)
	// eu = cceu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	// config.BlockNumber = new(big.Int).SetUint64(10000002)
	// config.Time = new(big.Int).SetUint64(10000002)

	// data = crypto.Keccak256([]byte("getSum()"))[:4]
	// msg = types.NewMessage(compiler.User1, &contractAddress, 2, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	// accesses, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{3, 3, 3}), 3, &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg))
	// t.Log("\n" + compiler.FormatTransitions(accesses))
	// t.Log("\n" + compiler.FormatTransitions(transitions))
	// t.Log(receipt)

	// if receipt.Status != 1 {
	// 	t.Error("Error: Set failed!!!", err)
	// }
}
