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
	"github.com/snippetor/bingo/node"
	"path/filepath"
	"github.com/snippetor/bingo/log"
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/app"
)

var (
	c          = flag.String("c", "", "config file path")
	d          = flag.Bool("d", false, "stop which node, with its name")
	n          = flag.String("n", "", "startup which node, with its name")
	proc       = flag.Int("p", -1, "cpu core size for runtime.GOMAXPROCS, default is runtime.NumCPU")
	help       = flag.Bool("h", false, "help")
	version    = flag.Bool("v", false, "bingo framework version")
	endRunning = make(chan bool, 1)
	usage      = `
	 Usage:
	 *if archive file name is echo
	 echo [options]
	 Or
	 go run echo.go [options]

	 Options:
	 -c : config file
	 -d : stop the node
	 -n : operate which node, with its name, if not set -n, all nodes will be operated
	 -p : cpu core size for runtime.GOMAXPROCS, default is runtime.NumCPU
	 -h : help
	 -v : bingo framework version

	 Example:
	 1. startup master node with config file bingo.json
	 echo -c bingo.json -n master
	 2. startup all nodes with config file bingo.json
	 echo -c bingo.json
	 3. stop master node
	 echo -c bingo.json -n master -d
	`
)

func Version() string {
	return "v1.0"
}

func BindNodeModel(modelName string, model interface{}) {
	app.BindNodeModel(modelName, model.(app.IModel))
}

func SetLogLevel(level log.Level) {
	fwlogger.SetLevel(level)
}

func Run() {
	flag.Parse()

	c := *c
	d := *d
	n := *n
	proc := *proc
	help := *help
	version := *version

	if help {
		fmt.Println(usage)
		return
	}
	if version {
		fmt.Println("Bingo version: " + Version())
		return
	}

	if c == "" {
		fmt.Println("config file must be set")
		fmt.Println(usage)
		return
	}

	if proc == -1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(proc)
	}

	if filepath.IsAbs(c) {
		app.Parse(c)
	} else {
		if dir, err := os.Getwd(); err == nil {
			app.Parse(filepath.Join(dir, c))
		}
	}

	if d {

	} else {
		if n == "" {
			app.RunAll()
		} else {
			app.Run(n)
		}
	}

	select {
	case <-endRunning:
		app.StopAll()
	}
	//<-endRunning
}
