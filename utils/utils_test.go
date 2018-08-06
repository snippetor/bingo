package utils

import (
	"testing"
	"fmt"
)

type A struct {
}

func TestStructName(t *testing.T) {
	a := A{}
	fmt.Println(StructName(a))
	fmt.Println(StructName(&a))
}

func TestElementName(t *testing.T) {
	var a []A
	fmt.Println(ElementName(a))
	fmt.Println(ElementName(&a))
}