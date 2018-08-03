package main

import (
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"strings"
	"bytes"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		dir, _ := os.Getwd()
		files, _ := ioutil.ReadDir(filepath.Join(dir, args[1]))
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			protofile := filepath.Join(dir, args[1], f.Name())
			if filepath.Ext(protofile) != ".proto" {
				continue
			}
			protoDir := filepath.Dir(protofile)
			genDir := filepath.Join(protoDir, strings.Split(f.Name(), ".")[0])
			os.Mkdir(genDir, os.ModePerm)
			cmd := exec.Command("tools/protoc", "--gogofaster_out="+genDir, "--proto_path="+protoDir, protofile)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(stderr.String())
			} else {
				fmt.Println("OK!")
			}
		}
	} else {
		fmt.Println("Usage: pb.exe ProtoFileDir")
	}
}
