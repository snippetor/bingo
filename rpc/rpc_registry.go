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
func RegisterMethod(methodName string, f RPCMethod) {
	if methods == nil {
		methods = make(map[string]RPCMethod)
	}
	methods[methodName] = f
}

func callMethod(methodName string, ctx *Context) *Result {
	if v, ok := methods[methodName]; ok {
		return v(ctx)
	}
	return nil
}
