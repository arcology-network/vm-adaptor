package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalizer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "runtime/runtime_test.sol", "0.8.19", "LocalizerTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
