package rpc

type ICallable interface {
	Call()
}

type IMethodRegistry interface {
	RegisterName(string, interface{})
	GetMethod(string) interface{}
}
