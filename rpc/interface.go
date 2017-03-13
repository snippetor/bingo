package rpc

type ICallable interface {
	Call()
	Go()
}

type IMethodRegistry interface {
	RegisterName(string, interface{})
	GetMethod(string) interface{}
}
