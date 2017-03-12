package net

import (
	"testing"
	"fmt"
)

func TestMessagePacker(t *testing.T) {
	p := &DefaultMessagePacker{}
	out := p.Pack(128, []byte("test_packer"))
	id, content, _ := p.Unpack(out)
	if id != 128 || string(content) != "test_packer" {
		fmt.Println(id, content)
		t.Fail()
	}
}

func TestMessagePacker1(t *testing.T) {
	p := &DefaultMessagePacker{}
	out := p.Pack(-123, []byte("中文测试"))
	id, content, _ := p.Unpack(out)
	if id != -123 || string(content) != "中文测试" {
		fmt.Println(id, content)
		t.Fail()
	}
}

func TestMessagePacker2(t *testing.T) {
	p := &DefaultMessagePacker{}
	out := p.Pack(-123, []byte(`{"a":"a", "b":1.1}`))
	id, content, _ := p.Unpack(out)
	if id != -123 || string(content) != `{"a":"a", "b":1.1}` {
		fmt.Println(id, content)
		t.Fail()
	}
}

func TestMessagePacker3(t *testing.T) {
	p := &DefaultMessagePacker{}
	out := p.Pack(-123, []byte(`{"a":"a", "b":1.1}`))
	out = append(out, []byte("[append]")...)
	id, content, out := p.Unpack(out)
	if id != -123 || string(content) != `{"a":"a", "b":1.1}` || string(out) != "[append]"{
		fmt.Println(id, content)
		t.Fail()
	}
}
