package main

import (
	"fmt"
	"reflect"
	"flag"
)

type A interface {
	func1()
}

type B struct {
}

func (b *B) func1() {
	fmt.Println("b")
}

func set(b *[]byte) {
	*b = []byte("test")
}

func main() {
	//dir, _ := os.Getwd()
	//fmt.Printf(dir)
	//f, _ := os.OpenFile(dir+"/log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	//defer f.Close()
	//l := log.New(f, "", log.Lmicroseconds)
	//l.Output(2, "[I] test")
	//l.Output(2, "[D] test1")
	var a A
	a = A(&B{})
	fmt.Println(reflect.TypeOf(a))
	(a).func1()

	addr := flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()
	fmt.Println(*addr)

	b := make([]byte, 0)
	set(&b)
	fmt.Println(string(b))
}
