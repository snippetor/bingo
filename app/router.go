package app

import (
	"fmt"
	"reflect"
	"github.com/snippetor/bingo/net"
	"strings"
)

type RouterBuilder interface {
	HandleWebApi(path interface{}, method string, middleware ...Handler)
	HandleRPCMethod(method string, middleware ...Handler)
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

func newRouterBuild(router Router, ctrl interface{}) RouterBuilder {
	return &routerBuilder{router: router, ctrl: ctrl}
}

func (r *routerBuilder) HandleWebApi(path interface{}, method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{Mix("API", path), method, middleware})
}

func (r *routerBuilder) HandleRPCMethod(method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{Mix("RPC", method), method, middleware})
}

func (r *routerBuilder) HandleServiceMessage(msgId interface{}, method string, middleware ...Handler) {
	r.methods = append(r.methods, &routerMethod{Mix("MSG", msgId), method, middleware})
}

func (r *routerBuilder) Build() {
	for _, m := range r.methods {
		if method, ok := reflect.TypeOf(r.ctrl).MethodByName(m.method); ok {
			var call = method.Func.Call
			r.router.Handle(m.key, func(ctx Context) {
				call([]reflect.Value{reflect.ValueOf(ctx)})
			})
		}
	}
}

type Router interface {
	Handle(key string, handlers ...Handler)
	Handlers(kind string) map[string]Handlers
	OnHandleRequest(ctx Context)
}

type router struct {
	routes map[string]Handlers
}

func NewRouter() Router {
	r := &router{make(map[string]Handlers)}
	return r
}

func (r *router) OnHandleRequest(ctx Context) {
	switch ctx.(type) {
	case *RpcContext:
		mixed := Mix("RPC", ctx.(*RpcContext).Method)
		if hs, ok := r.routes[mixed]; ok {
			var newHandlers Handlers
			globalMiddleWares := ctx.App().GlobalMiddleWares()
			if r.apply(&globalMiddleWares) {
				newHandlers = append(newHandlers, globalMiddleWares...)
			}
			newHandlers = append(newHandlers, hs...)
			ctx.Do(newHandlers)
		}
	case *WebApiContext:
		mixed := Mix("API", string(ctx.(*WebApiContext).RequestCtx.Path()))
		if hs, ok := r.routes[mixed]; ok {
			var newHandlers Handlers
			globalMiddleWares := ctx.App().GlobalMiddleWares()
			if r.apply(&globalMiddleWares) {
				newHandlers = append(newHandlers, globalMiddleWares...)
			}
			newHandlers = append(newHandlers, hs...)
			ctx.Do(newHandlers)
		}
	case *ServiceContext:
		c := ctx.(*ServiceContext)
		if c.MessageType != net.MsgTypeReq {
			c.LogE("Ignore service message type=%d, group=%d, extra=%d id=%d", c.MessageType, c.MessageBody, c.MessageExtra, c.MessageId)
			return
		}
		var mixed = []string{
			Mix("MSG", net.PackId(c.MessageType, c.MessageGroup, c.MessageExtra, c.MessageId)),
			Mix("MSG", net.PackId(c.MessageType, c.MessageGroup, c.MessageExtra, 0)),
			Mix("MSG", net.PackId(c.MessageType, c.MessageGroup, 0, 0)),
		}
		for _, m := range mixed {
			if hs, ok := r.routes[m]; ok {
				var newHandlers Handlers
				globalMiddleWares := c.App().GlobalMiddleWares()
				if r.apply(&globalMiddleWares) {
					newHandlers = append(newHandlers, globalMiddleWares...)
				}
				newHandlers = append(newHandlers, hs...)
				c.Do(newHandlers)
			}
		}
	}

}

func (r *router) Handle(key string, handlers ...Handler) {
	mainHandlers := Handlers(handlers)
	if !r.apply(&mainHandlers) {
		return
		if hs, ok := r.routes[key]; ok {
			r.routes[key] = append(hs, mainHandlers...)
		} else {
			r.routes[key] = mainHandlers
		}
	}
}

func (r *router) Handlers(kind string) map[string]Handlers {
	m := make(map[string]Handlers)
	for key, handlers := range r.routes {
		if strings.HasPrefix(key, kind) {
			m[key] = handlers
		}
	}
	return m
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

func Mix(method string, path interface{}) string {
	return fmt.Sprintf("%s:%v", method, path)
}
