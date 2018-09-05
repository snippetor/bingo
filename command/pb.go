package command

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"os/exec"
	"strings"
	"github.com/snippetor/bingo/utils"
	"runtime"
)

func Pb(srcDir, outputDir string, split bool) {
	dir, _ := os.Getwd()
	files, _ := ioutil.ReadDir(filepath.Join(dir, srcDir))
	if len(files) == 0 {
		printError("Not found .proto file in %s", srcDir)
		return
	}
	protoc, ok := checkExecFiles()
	if !ok {
		return
	}
	var hasProto bool
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		protofile := filepath.Join(dir, srcDir, f.Name())
		if filepath.Ext(protofile) != ".proto" {
			continue
		}
		hasProto = true
		protoDir := filepath.Dir(protofile)
		if outputDir == "" {
			outputDir = protoDir
		}
		var cmd *exec.Cmd
		if split {
			genDir := filepath.Join(outputDir, strings.Split(f.Name(), ".")[0])
			os.Mkdir(genDir, os.ModePerm)
			cmd = exec.Command(protoc, "--gogofaster_out="+genDir, "--proto_path="+protoDir, protofile)
		} else {
			cmd = exec.Command(protoc, "--gogofaster_out="+outputDir, "--proto_path="+protoDir, protofile)
		}
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(err)
		} else {
			printSuccess("pb: %s OK", f.Name())
		}
	}
	if !hasProto {
		printError("Not found .proto file in %s", srcDir)
	}
}

func checkExecFiles() (string, bool) {
	var rootDir []string
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		rootDir = append(rootDir, gopath)
	}
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		rootDir = append(rootDir, goroot)
	}

	if len(rootDir) > 0 {
		var protoc string
		var ok bool
		for i := range rootDir {
			if runtime.GOOS == "windows" {
				protoc = filepath.Join(rootDir[i], "src", "github.com", "snippetor", "bingo", "tools", "protoc.exe")
			} else {
				protoc = filepath.Join(rootDir[i], "src", "github.com", "snippetor", "bingo", "tools", "protoc")
			}
			if utils.IsFileExists(protoc) {
				ok = true
				break
			}
		}
		if !ok {
			printError("Not found protoc in $GOPATH/src/github.com/snippetor/bingo/tools/")
			return "", false
		}
		ok = false
		var gogofaster string
		for i := range rootDir {
			if runtime.GOOS == "windows" {
				gogofaster = filepath.Join(rootDir[i], "bin", "protoc-gen-gogofaster.exe")
			} else {
				gogofaster = filepath.Join(rootDir[i], "bin", "protoc-gen-gogofaster")
			}
			if utils.IsFileExists(gogofaster) {
				ok = true
				break
			}
		}
		if !ok {
			printError("Not found protoc-gen-gogofaster, please exec 'go get github.com/gogo/protobuf/protoc-gen-gogofaster'")
			return "", false
		}
		return protoc, true
	} else {
		printError("Not found GOPATH or GOROOT")
		return "", false
	}
}
