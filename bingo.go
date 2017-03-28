package bingo

import (
	"github.com/snippetor/bingo/log"
	"fmt"
	"flag"
	"runtime"
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
	config = flag.String("c", "", "config file path or file name in current fold")
	node   = flag.String("n", "", "startup which node, with its name")
	all    = flag.Bool("a", false, "startup all nodes in one server")
)

func Run() {
	// Usage:
	// *if archive file name is echo
	// echo [options]
	//
	// Options:
	// -c : config file path or file name in current fold
	// -n : startup which node, with its name
	// -a : startup all nodes in one server
	// -p : cpu core size for runtime.GOMAXPROCS, default is runtime.NumCPU
	// -h : help
	//
	// Example:
	// 1. startup master node with config file bingo.json
	// echo -c bingo.json -n master
	// 2. startup all nodes with config file bingo.json
	// echo -c bingo.json -a
	flag.Parse()

	config := *config
	node := *node
	all := *all

	if config == "" {
		fmt.Println("-c must be set with config file path or file name in current fold")
		return
	}
	if !all && node == "" {
		fmt.Println("-n must be set when not use -a for startup all nodes")
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(config)
	fmt.Println(node)
	fmt.Println(all)

}
