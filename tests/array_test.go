package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestAddressContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/address_test.sol", "0.8.19", "AddressTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBoolContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/bool_test.sol", "0.8.19", "BoolTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBytesContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/bytes_test.sol", "0.8.19", "ByteTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractBytes32(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/bytes32_test.sol", "0.8.19", "Bytes32Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/u256_test.sol", "0.8.19", "U256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractInt256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/int256_test.sol", "0.8.19", "Int256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractString(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _, _ := DeployThenInvoke(targetPath, "array/string_test.sol", "0.8.19", "StringTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _, _ := DeployThenInvoke(targetPath, "commutative/u256Cum_test.sol", "0.8.19", "CumulativeU256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestU256Multiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _, _ := DeployThenInvoke(targetPath, "/multiprocess/mp_u256_test.sol", "0.8.19", "U256ParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestArrayMultiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _, _ := DeployThenInvoke(targetPath, "/multiprocess/mp_u256_test.sol", "0.8.19", "U256ArrayParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParallelArrayPush(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _, _ := DeployThenInvoke(targetPath, "/multiprocess/multiprocess_test.sol", "0.8.19", "ParaNativeArrayPushbackTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
