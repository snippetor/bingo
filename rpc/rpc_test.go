package rpc

import (
	"testing"
	"fmt"
	"github.com/snippetor/bingo/codec"
	"unsafe"
)

func TestGenCallSeq(t *testing.T) {
}

func TestMap(t *testing.T) {
	c := &RPCMethodCall{123, "aaa.aa", nil, ""}
	cd := codec.NewCodec(codec.Protobuf)
	fmt.Println(c.Args == nil)
	bs := cd.Marshal(c)
	c1 := &RPCMethodCall{}
	cd.Unmarshal(bs, c1)
	fmt.Println(c1)
	fmt.Println(c1.Args == nil)

	v := RPCValue{I32: 0, S: "test"}
	fmt.Println(unsafe.Sizeof(v))
	bytes := make([]byte, v.Size())
	v.MarshalTo(bytes)
	fmt.Println(len(bytes))
}
