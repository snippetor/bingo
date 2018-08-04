package route

import (
	"reflect"
	"runtime"
)

type Handler func(Context)

// Handlers is just a type of slice of []Handler.
//
// See `Handler` for more.
type Handlers []Handler

// HandlerName returns the name, the handler function informations.
// Same as `context.HandlerName`.
func HandlerName(h Handler) string {
	pc := reflect.ValueOf(h).Pointer()
	return runtime.FuncForPC(pc).Name()
}
