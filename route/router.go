package route

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type Router interface {
	HandleWebApi(path interface{}, handlers ...Handler)
	HandleRPCMethod(method string, handlers ...Handler)
	HandleServiceMessage(msgId interface{}, handlers ...Handler)
}

var _ Router = (*router)(nil)

type router struct {
	routes map[interface{}]Handlers
}

func NewRouter() Router {
	r := &router{make(map[interface{}]Handlers)}
	return r
}

func (r *router) HandleWebApi(path interface{}, handlers ...Handler) {
	mainHandlers := Handlers(handlers)
	if !r.apply(&mainHandlers) {
		return
	}
	mixed := r.mix("API", path)
	if hs, ok := r.routes[mixed]; ok {
		r.routes[mixed] = append(hs, mainHandlers...)
	} else {
		r.routes[mixed] = mainHandlers
	}
}

func (r *router) HandleRPCMethod(method string, handlers ...Handler) {
	mainHandlers := Handlers(handlers)
	if !r.apply(&mainHandlers) {
		return
	}
	mixed := r.mix("RPC", method)
	if hs, ok := r.routes[mixed]; ok {
		r.routes[mixed] = append(hs, mainHandlers...)
	} else {
		r.routes[mixed] = mainHandlers
	}
}

func (r *router) HandleServiceMessage(msgId interface{}, handlers ...Handler) {
	mainHandlers := Handlers(handlers)
	if !r.apply(&mainHandlers) {
		return
	}
	mixed := r.mix("MSG", msgId)
	if hs, ok := r.routes[mixed]; ok {
		r.routes[mixed] = append(hs, mainHandlers...)
	} else {
		r.routes[mixed] = mainHandlers
	}
}

func (r *router) OnWebApiRequest(ctx *fasthttp.RequestCtx) {
	mixed := r.mix("API", string(ctx.Path()))
	if hs, ok := r.routes[mixed]; ok {
		//ctx.Do(hs)
	}
}

func (r *router) OnRPCRequest(method string, ctx Context) {
	mixed := r.mix("RPC", method)
	if hs, ok := r.routes[mixed]; ok {
		ctx.Do(hs)
	}
}

func (r *router) OnServiceRequest(msgId interface{}, ctx Context) {
	mixed := r.mix("MSG", msgId)
	if hs, ok := r.routes[mixed]; ok {
		ctx.Do(hs)
	}
}

func (r router) buildHandler(h Handler) Handler {
	return func(ctx Context) {
		// Proceed will fire the handler and return false here if it doesn't contain a `ctx.Next()`,
		// so we add the `ctx.Next()` wherever is necessary in order to eliminate any dev's misuse.
		if !ctx.Proceed(h) {
			// `ctx.Next()` always checks for `ctx.IsStopped()` and handler(s) positions by-design.
			ctx.Next()
		}
	}
}

func (r router) apply(handlers *Handlers) bool {
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
