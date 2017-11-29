package utils

import (
	"testing"
	"fmt"
	"reflect"
)

func TestValue_Set(t *testing.T) {
	v := Value{}
	v.Set(&Value{})
	fmt.Println(reflect.TypeOf(v.inner).Kind())
}