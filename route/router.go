package route

import (
	"fmt"
	"github.com/snippetor/bingo/net"
	"reflect"
)

type RouterBuilder interface {
	HandleWebApi(path interface{}, method string, middleware ...Handler)
	HandleRPCMethod(method string, handlers ...Handler)
	HandleServiceMessage(msgId interface{}, method string, middleware ...Handler)
	Build()
}

type routerMethod struct {
	key        string
	method     string
	middleware Handlers
}

type routerBuilder struct {
	router  Router
	ctrl    interface{}
	methods []*routerMethod
}

func NewRouterBuild(router Router, ctrl interface{}) RouterBuilder {
	return &routerBuilder{router: router, ctrl: ctrl}
}

func (r *routerBuilder) HandleWebApi(path interface{}, method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{r.router.mix("API", path), method, middleware})
}

func (r *routerBuilder) HandleRPCMethod(method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{r.router.mix("RPC", method), method, middleware})
}

func (r *routerBuilder) HandleServiceMessage(msgId interface{}, method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{r.router.mix("MSG", msgId), method, middleware})
}

func (r *routerBuilder) Build() {
	for _, m := range r.methods {
		if method, ok := reflect.TypeOf(r.ctrl).MethodByName(m.method); ok {
			var call = method.Func.Call
			r.router.handle(m.key, func(ctx Context) {
				call([]reflect.Value{reflect.ValueOf(ctx)})
			})
		}
	}
}

type Router interface {
	OnWebApiRequest(path string, ctx Context)
	OnRPCRequest(method string, ctx Context)
	OnServiceRequest(msgId net.MessageId, ctx Context)

	handle(key string, handlers ...Handler)
	mix(method string, path interface{}) string
}

type router struct {
	routes map[interface{}]Handlers
}

func NewRouter() Router {
	r := &router{make(map[interface{}]Handlers)}
	return r
}

func (r *router) OnWebApiRequest(path string, ctx Context) {
	mixed := r.mix("API", path)
	if hs, ok := r.routes[mixed]; ok {
		var newHandlers Handlers
		globalMiddleWares := ctx.App().GlobalMiddleWares()
		if r.apply(&globalMiddleWares) {
			newHandlers = append(newHandlers, globalMiddleWares...)
		}
		newHandlers = append(newHandlers, hs...)
		ctx.Do(newHandlers)
	}
}

func (r *router) OnRPCRequest(method string, ctx Context) {
	mixed := r.mix("RPC", method)
	if hs, ok := r.routes[mixed]; ok {
		var newHandlers Handlers
		globalMiddleWares := ctx.App().GlobalMiddleWares()
		if r.apply(&globalMiddleWares) {
			newHandlers = append(newHandlers, globalMiddleWares...)
		}
		newHandlers = append(newHandlers, hs...)
		ctx.Do(newHandlers)
	}
}

func (r *router) OnServiceRequest(msgId net.MessageId, ctx Context) {
	mixed := r.mix("MSG", msgId)
	if hs, ok := r.routes[mixed]; ok {
		var newHandlers Handlers
		globalMiddleWares := ctx.App().GlobalMiddleWares()
		if r.apply(&globalMiddleWares) {
			newHandlers = append(newHandlers, globalMiddleWares...)
		}
		newHandlers = append(newHandlers, hs...)
		ctx.Do(newHandlers)
	}
}

func (r *router) handle(key string, handlers ...Handler) {
	mainHandlers := Handlers(handlers)
	if !r.apply(&mainHandlers) {
		return
	}
	if hs, ok := r.routes[key]; ok {
		r.routes[key] = append(hs, mainHandlers...)
	} else {
		r.routes[key] = mainHandlers
	}
}

func (r *router) buildHandler(h Handler) Handler {
	return func(ctx Context) {
		// Proceed will fire the handler and return false here if it doesn't contain a `ctx.Next()`,
		// so we add the `ctx.Next()` wherever is necessary in order to eliminate any dev's misuse.
		if !ctx.Proceed(h) {
			// `ctx.Next()` always checks for `ctx.IsStopped()` and handler(s) positions by-design.
			ctx.Next()
		}
	}
}

func (r *router) apply(handlers *Handlers) bool {
	tmp := *handlers
	for i, h := range tmp {
		if h == nil {
			if len(tmp) == 1 {
				return false
			}
			continue
		}
		(*handlers)[i] = r.buildHandler(h)
	}
	return true
}

func (r *router) mix(method string, path interface{}) string {
	return fmt.Sprintf("%s:%v", method, path)
}
