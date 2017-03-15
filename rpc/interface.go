package rpc

type ICallable interface {
	Call(string, map[string]string, string) bool
}

type IMethodRegistry interface {
	RegisterName(string, interface{})
	GetMethod(string) interface{}
}
