package utils

import (
	"testing"
	"fmt"
)

func TestValue_Set(t *testing.T) {
	v := Value{}
	v.Set("test")
	fmt.Println(v.GetString())


	vm := &ValueMap{}
	vm.Put("test", "test")
	fmt.Println(vm.Get("test").GetString())
}