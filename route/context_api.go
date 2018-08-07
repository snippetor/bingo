package route

import (
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/codec"
)

type WebApiContext struct {
	Context
	RequestCtx *fasthttp.RequestCtx
	Codec      codec.ICodec
}

// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*MyContext".
func (c *WebApiContext) Do(handlers Handlers) {
	Do(c, handlers)
}

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*MyContext".
func (c *WebApiContext) Next() {
	Next(c)
}

func (c *WebApiContext) RequestBody(body interface{}) {
	c.Codec.Unmarshal(c.RequestCtx.Request.Body(), body)
}

func (c *WebApiContext) ResponseOK(body interface{}) {
	c.RequestCtx.Response.SetStatusCode(fasthttp.StatusOK)
	bs := c.Codec.Marshal(body)
	c.RequestCtx.Response.SetBody(bs)
	c.LogD("<<< %s %s", string(c.RequestCtx.Path()), string(bs))
}

func (c *WebApiContext) ResponseFailed(reason string) {
	c.RequestCtx.Response.SetStatusCode(fasthttp.StatusOK)
	params := make(map[string]interface{})
	params["error"] = reason
	bs := c.Codec.Marshal(params)
	c.RequestCtx.Response.SetBody(bs)
	c.LogD("<<< %s %s", string(c.RequestCtx.Path()), string(bs))
}
