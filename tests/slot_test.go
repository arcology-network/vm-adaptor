package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	common "github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	ccEu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
	"github.com/holiman/uint256"
	sha3 "golang.org/x/crypto/sha3"
)

func TestSlotHash(t *testing.T) {
	_ctrn := uint256.NewInt(0).Bytes32()
	_elem := uint256.NewInt(2).Bytes32()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(append(_elem[:], _ctrn[:]...))
	fmt.Println(hash.Sum(nil))
}

func TestStorageSlot(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/int256/Base.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/bool/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, project+"/api/slot/local_test.sol", "LocalTest")
	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

}