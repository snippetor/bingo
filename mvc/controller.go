package mvc

import (
	"github.com/snippetor/bingo/app"
	"reflect"
)

var controllerType = reflect.TypeOf((*Controller)(nil)).Elem()

func IsController(i interface{}) bool {
	return reflect.TypeOf(i).Implements(controllerType)
}

type Controller interface {
	Route(builder app.RouterBuilder)
}
