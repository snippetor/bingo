package mvc

import "reflect"

var controllerType = reflect.TypeOf((*Controller)(nil)).Elem()
func IsController(i interface{}) bool {
	return reflect.TypeOf(i).Implements(controllerType)
}

var ormModelType = reflect.TypeOf((*OrmModel)(nil)).Elem()
func IsOrmModel(i interface{}) bool {
	return reflect.TypeOf(i).Implements(ormModelType)
}
