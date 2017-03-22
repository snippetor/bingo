package rpc

import (
	"reflect"
	"github.com/snippetor/bingo"
)

var (
	methods map[string]reflect.Value
)

// 注册目标对象的所有方法
// v：必须是指针
func RegisterMethods(target string, v interface{}) {
	if methods == nil {
		methods = make(map[string]reflect.Value)
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		bingo.E("-- register methods failed! v must be a pointer. --")
		return
	}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		methods[makeKey(target, t.Method(i).Name)] = reflect.ValueOf(v).Method(i)
	}
}

// 注册单个方法
func RegisterMethod(target, method string, v interface{}) {
	if methods == nil {
		methods = make(map[string]reflect.Value)
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		bingo.E("-- register method failed! v must be a pointer. --")
		return
	}
	methods[makeKey(target, method)] = reflect.ValueOf(v).MethodByName(method)
}

func callMethod(target, method string, ctx *Context) {
	if v, ok := methods[makeKey(target, method)]; ok {
		v.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

func makeKey(target, method string) string {
	return target + "." + method
}
