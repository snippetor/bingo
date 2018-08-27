package recover

import (
	"fmt"
	"runtime"
	"reflect"
	"github.com/snippetor/bingo/app"
)

func getRequestLogs(ctx app.Context) string {
	switch ctx.(type) {
	case *app.RpcContext:
		c := ctx.(*app.RpcContext)
		return fmt.Sprintf("[RPC] call '%s' method %s", c.App().Name(), c.Method)
	case *app.ServiceContext:
		c := ctx.(*app.ServiceContext)
		return fmt.Sprintf("[SOC] '%s' %v %v %v %v, %v", c.App().Name(), c.MessageType, c.MessageGroup, c.MessageExtra, c.MessageId, c.MessageBody.RawContent)
	case *app.WebApiContext:
		c := ctx.(*app.WebApiContext)
		return fmt.Sprintf("[API] '%s' %s %s %s", c.App().Name(), string(c.RequestCtx.Path()), string(c.RequestCtx.Method()), c.RequestCtx.RemoteIP().String())
	default:
		return fmt.Sprintf("[UFO] unknown ctx type %v", reflect.TypeOf(ctx))
	}
}

// New returns a new recover middleware,
// it recovers from panics and logs
// the panic message to the application's logger "Warn" level.
func New() app.Handler {
	return func(ctx app.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}

				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break

					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)
				ctx.LogE(logMessage)
				switch ctx.(type) {
				case *app.RpcContext:
					c := ctx.(*app.RpcContext)
					c.ReturnNil()
					//case *app.ServiceContext:
					//c := ctx.(*app.ServiceContext)
					//c.Ack(map[string]interface{}{"code": -1, "desc": fmt.Sprintf("%s", err)})
				case *app.WebApiContext:
					c := ctx.(*app.WebApiContext)
					c.ResponseFailed(fmt.Sprintf("%s", err))
				}
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}
