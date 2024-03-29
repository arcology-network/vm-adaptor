package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/arcology-network/common-lib/common"
)

func GetContractMeta(file string) (contractName string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
		if strings.HasPrefix(line, "contract") {
			words := strings.Split(line, " ")
			contractName = words[1]
			break
		}
	}
	return
}

func CompileContracts(dockerRootpath, solfilename, version, contractname string, outpathhold bool) (string, error) {
	if !outpathhold {
		removeOut(dockerRootpath)
	}

	if !common.FileExists(dockerRootpath + "/" + solfilename) {
		return "", errors.New("Error: The contract file doesn't exist in " + dockerRootpath + "/" + solfilename)
	}

	ensureOutpath(dockerRootpath)

	if _, err := exec.Command(
		"docker", "run",
		"-v", dockerRootpath+":/sources",
		"ethereum/solc:"+version,
		"-o", "/sources/"+outpath,
		"--abi", "--bin", "--overwrite",
		"/sources/"+solfilename).Output(); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(dockerRootpath + "/" + outpath + "/" + contractname + ".bin")
	if err != nil {
		fmt.Printf("reading contract err:%v\n", err)
		return "", err
	}
	if !outpathhold {
		removeOut(dockerRootpath)
	}
	return fmt.Sprintf("%s", bytes), nil
}

const (
	outpath = "output"
)

func ensureOutpath(basePath string) {
	folderPath := filepath.Join(basePath, outpath)
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	} else {
		dir, _ := ioutil.ReadDir(folderPath)
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{folderPath, d.Name()}...))
		}
	}
}

func removeOut(basePath string) {
	folderPath := filepath.Join(basePath, outpath)
	dir, _ := ioutil.ReadDir(folderPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{folderPath, d.Name()}...))
	}
	os.Remove(folderPath)
}
