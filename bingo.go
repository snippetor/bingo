// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bingo

import (
	"fmt"
	"flag"
	"runtime"
	"os"
	"sort"
	"strings"
	"github.com/snippetor/bingo/node"
	"path/filepath"
)

var (
	n          = flag.String("n", "", "startup which node, with its name")
	proc       = flag.Int("p", -1, "cpu core size for runtime.GOMAXPROCS, default is runtime.NumCPU")
	help       = flag.Bool("h", false, "help")
	version    = flag.Bool("v", false, "bingo framework version")
	endRunning = make(chan bool, 1)
	usage      = `
	 Usage:
	 *if archive file name is echo
	 echo [command] [config file] [options]

	 Command:
	 start : startup node
	 stop  : stop node

	 Options:
	 -n : operate which node, with its name, if not set -n, all nodes will be operated
	 -p : cpu core size for runtime.GOMAXPROCS, default is runtime.NumCPU
	 -h : help
	 -v : bingo framework version

	 Example:
	 1. startup master node with config file bingo.json
	 echo -c bingo.json -n master
	 2. startup all nodes with config file bingo.json
	 echo start bingo.json -a or echo start bingo.json
	`
	commands = []string{"start", "stop"}
)

func Version() string {
	return "1.0"
}

func BindNodeModel(modelName string, model interface{}) {
	node.BindNodeModel(modelName, model.(node.IModel))
}

func Run() {

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}

	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}
	cmd := strings.ToLower(os.Args[1])
	config := os.Args[2]

	if sort.SearchStrings(commands, cmd) < 0 {
		fmt.Println("command must be one of", commands)
		return
	}

	if config == "" {
		fmt.Println("config file must be set")
		return
	}

	flag.Parse()

	n := *n
	proc := *proc
	help := *help
	version := *version

	if help {
		fmt.Println(usage)
		return
	}
	if version {
		fmt.Println(Version())
		return
	}

	if proc == -1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(proc)
	}

	if filepath.IsAbs(config) {
		node.Parse(config)
	} else {
		if dir, err := os.Getwd(); err == nil {
			node.Parse(filepath.Join(dir, config))
		}
	}

	switch cmd {
	case "start":
		if n == "" {
			node.RunAll()
		} else {
			node.Run(n)
		}
	case "stop":

	}

	<-endRunning
}
