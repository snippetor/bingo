package node

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

type httpContext struct {
}

type httpServer struct {
}

func (s *httpServer) Listen(port int) {
	f := func(ctx *fasthttp.RequestCtx) {
	}
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), f)
}