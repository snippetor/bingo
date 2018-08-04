package route

import (
	"github.com/valyala/fasthttp"
	"github.com/bitly/go-simplejson"
	"github.com/snippetor/bingo/errors"
	"encoding/json"
)

type WebApiContext struct {
	Context
	RequestCtx *fasthttp.RequestCtx
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

func (c *WebApiContext) RequestBody() *simplejson.Json {
	j, err := simplejson.NewJson(c.RequestCtx.Request.Body())
	errors.Check(err)
	return j
}

func (c *WebApiContext) ResponseOK(body interface{}) {
	c.RequestCtx.Response.SetStatusCode(fasthttp.StatusOK)
	if bs, err := json.Marshal(body); err == nil {
		c.RequestCtx.Response.SetBody(bs)
		c.LogD("<==== %s %s", string(c.RequestCtx.Path()), string(bs))
	} else {
		c.LogE("send Failed !!!! %s", err.Error())
	}
}

func (c *WebApiContext) ResponseFailed(reason string) {
	c.RequestCtx.Response.SetStatusCode(fasthttp.StatusOK)
	j := simplejson.New()
	j.Set("error", reason)
	if bs, err := j.Encode(); err == nil {
		c.RequestCtx.Response.SetBody(bs)
		c.LogD("<==== %s %s", string(c.RequestCtx.Path()), string(bs))
	} else {
		c.LogE("send Failed !!!! %s", err.Error())
	}
}
