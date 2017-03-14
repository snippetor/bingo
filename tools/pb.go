package main

import (
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		protofile := args[1]
		if filepath.Ext(protofile) != ".proto" {
			fmt.Println("must be .proto file")
			return
		}
		dir := filepath.Dir(protofile)
		err := exec.Command("protoc", "--gogofaster_out="+dir, "--proto_path="+dir, protofile).Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Usage: pb.exe ProtoFilePath")
	}
}
