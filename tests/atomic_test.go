package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	concurrenturl "github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/crypto"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
)

func TestAtomicWithThreading(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/"

	// if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
	// 	t.Error(err)
	// }

	code, err := compiler.CompileContracts(targetPath, "atomic/atomic_test.sol", "0.8.19", "AtomicDeferredInThreadingTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, project+"/api/atomic/atomic_test.sol", "AtomicDeferredInThreadingTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	contractAddress := receipt.ContractAddress
	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}

	data = crypto.Keccak256([]byte("PostCheck()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}

func TestAtomicWithThreadingAndContainer(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/"

	// if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
	// 	t.Error(err)
	// }
	code, err := compiler.CompileContracts(targetPath, "atomic/atomic_test.sol", "0.8.19", "AtomicDeferredBoolContainerTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, project+"/api/atomic/atomic_test.sol", "AtomicDeferredBoolContainerTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	contractAddress := receipt.ContractAddress
	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}

	data = crypto.Keccak256([]byte("PostCheck()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}

func TestAtomicMultiDeferredWithBoolContainer(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/"

	// if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
	// 	t.Error(err)
	// }
	code, err := compiler.CompileContracts(targetPath, "atomic/atomic_test.sol", "0.8.19", "AtomicMultiDeferredWithBoolContainerTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, project+"/api/atomic/atomic_test.sol", "AtomicMultiDeferredWithBoolContainerTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	contractAddress := receipt.ContractAddress
	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}

	data = crypto.Keccak256([]byte("PostCheck()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}

func TestAtomicMultiDeferreOneConflictTest(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/"

	// if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
	// 	t.Error(err)
	// }
	code, err := compiler.CompileContracts(targetPath, "atomic/atomic_test.sol", "0.8.19", "AtomicMultiDeferreOneConflictTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, project+"/api/atomic/atomic_test.sol", "AtomicMultiDeferreOneConflictTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	if err != nil {
		fmt.Print(err)
	}

	contractAddress := receipt.ContractAddress
	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}
