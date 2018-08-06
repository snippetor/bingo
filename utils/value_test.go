package utils

import (
	"testing"
	"fmt"
)

func TestValue_Set(t *testing.T) {
	v := NewValue()
	v.Set("test")
	fmt.Println(v.GetString())


	vm := NewValueMap()
	vm.Put("test", "test")
	fmt.Println(vm.Get("test").GetString())
}