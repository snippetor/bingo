package rpc

type Result struct {
	Args
}

type RPCMethod func(*Context) *Result
type RPCCallback func(*Result)

var (
	methods map[string]RPCMethod
)

// 注册单个方法
func RegisterMethod(target, methodName string, f RPCMethod) {
	if methods == nil {
		methods = make(map[string]RPCMethod)
	}
	methods[makeKey(target, methodName)] = f
}

func callMethod(target, methodName string, ctx *Context) *Result {
	if v, ok := methods[makeKey(target, methodName)]; ok {
		return v(ctx)
	}
	return nil
}

func makeKey(target, methodName string) string {
	return target + "." + methodName
}
