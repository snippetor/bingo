package main

import (
	"fmt"
	"reflect"
)

type A interface {
	func1()
}

type B struct {
}

func (b *B) func1() {
	fmt.Println("b")
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
}
