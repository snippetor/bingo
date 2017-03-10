package net

import (
	"testing"
	"fmt"
)

type JsonObj struct {
	Id   int
	Name string
}

func TestProtocolJson(t *testing.T) {
	p := ProtocolJson{}
	bs, err := p.marshal(&JsonObj{1, "carl"})
	if err != nil {
		t.Fail()
	}

	var b JsonObj
	err = p.unmarshal(bs, &b)
	if err != nil {
		t.Fail()
	}
	fmt.Println(b)
}

func TestProtocolProtobuf(t *testing.T) {
	p := ProtocolProtoBuf{}
	persion := &Person{}
	var id int32 = 1
	var name string = "carl"
	persion.Id = &id
	persion.Name = &name
	bs, err := p.marshal(persion)
	if err != nil {
		t.Fail()
	}

	persion2 := Person{}
	err = p.unmarshal(bs, &persion2)
	if err != nil {
		t.Fail()
	}
	fmt.Println(*persion2.Name)
}