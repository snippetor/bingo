package recover

import (
	"fmt"
	"runtime"
	"github.com/snippetor/bingo/route"
	"reflect"
)

func getRequestLogs(ctx route.Context) string {
	switch ctx.(type) {
	case *route.RpcContext:
		c := ctx.(*route.RpcContext)
		return fmt.Sprintf("[RPC] '%s' call '%s' method %s#%v args=%v", c.Caller, c.App().Name(), c.Method, c.CallSeq, c.Args)
	case *route.ServiceContext:
		c := ctx.(*route.ServiceContext)
		return fmt.Sprintf("[SOC] '%s' %v %v %v %v, %v", c.App().Name(), c.MessageType, c.MessageGroup, c.MessageExtra, c.MessageId, c.MessageBody.RawContent)
	case *route.WebApiContext:
		c := ctx.(*route.WebApiContext)
		return fmt.Sprintf("[API] '%s' %s %s %s", c.App().Name(), string(c.RequestCtx.Path()), string(c.RequestCtx.Method()), c.RequestCtx.RemoteIP().String())
	default:
		return fmt.Sprintf("[UFO] unknown ctx type %v", reflect.TypeOf(ctx))
	}
}

// New returns a new recover middleware,
// it recovers from panics and logs
// the panic message to the application's logger "Warn" level.
func New() route.Handler {
	return func(ctx route.Context) {
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
				case *route.RpcContext:
					c := ctx.(*route.RpcContext)
					c.ReturnNil()
					//case *route.ServiceContext:
					//c := ctx.(*route.ServiceContext)
					//c.Ack(map[string]interface{}{"code": -1, "desc": fmt.Sprintf("%s", err)})
				case *route.WebApiContext:
					c := ctx.(*route.WebApiContext)
					c.ResponseFailed(fmt.Sprintf("%s", err))
				}
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}
