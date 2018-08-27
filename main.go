package main

import (
	"reflect"
	"runtime"
	"fmt"
)

type A struct {
}

func Test() {
}

func HandlerName(h func()) string {
	pc := reflect.ValueOf(h).Pointer()
	return runtime.FuncForPC(pc).Name()
}
func main() {

	fmt.Println(HandlerName(Test))
	//flag.Parse()
	//
	//c := *c
	//d := *d
	//n := *n
	//proc := *proc
	//help := *help
	//version := *version
	//
	//if help {
	//	fmt.Println(usage)
	//	return
	//}
	//if version {
	//	fmt.Println("Bingo version: " + Version())
	//	return
	//}
	//
	//if c == "" {
	//	fmt.Println("config file must be set")
	//	fmt.Println(usage)
	//	return
	//}
	//
	//if proc == -1 {
	//	runtime.GOMAXPROCS(runtime.NumCPU())
	//} else {
	//	runtime.GOMAXPROCS(proc)
	//}

	//if filepath.IsAbs(c) {
	//	parse(c)
	//} else {
	//	if dir, err := os.Getwd(); err == nil {
	//		parse(filepath.Join(dir, c))
	//	}
	//}
	//
	//if d {
	//	//TODO
	//} else {
	//	if n == "" {
	//		runAll()
	//	} else {
	//		runApp(n)
	//	}
	//}
	//
	//select {
	//case <-endRunning:
	//	stopAll()
	//}
	//<-endRunning
}
