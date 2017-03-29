package bingo

import (
	"github.com/snippetor/bingo/log"
	"fmt"
	"flag"
	"runtime"
	"os"
	"sort"
)

var (
	// bingo框架日志
	fwLogger *log.Logger
)

func init() {
	fwLogger = log.NewLoggerWithConfig(log.DEFAULT_CONFIG)
}

func Version() string {
	return "1.0"
}

func I(format string, v ...interface{}) {
	fwLogger.I(format, v...)
}

func D(format string, v ...interface{}) {
	fwLogger.D(format, v...)
}

func W(format string, v ...interface{}) {
	fwLogger.W(format, v...)
}

func E(format string, v ...interface{}) {
	fwLogger.E(format, v...)
}

var (
	node       = flag.String("n", "", "startup which node, with its name")
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

func Run() {

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}

	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}
	cmd := os.Args[1]
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

	node := *node
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

	<-endRunning
}
