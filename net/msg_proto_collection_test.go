package net

import (
	"testing"
	"fmt"
)

type TestObj struct {
	Id   int
	Name string
}

func TestProtoCollection(t *testing.T) {
	c := protoCollection{}
	c.put(1, TestObj{})

	c1, ok := c.get(1)
	if ok {
		o := c1.(TestObj)
		o.Id = 1
	}

	c2, ok := c.get(1)
	if ok {
		o := c2.(TestObj)
		fmt.Println(o.Id)
	}
}
