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
	if bs, err := cd.Marshal(c); err == nil {
		c1 := &RPCMethodCall{}
		if err := cd.Unmarshal(bs, c1); err == nil {
			fmt.Println(c1)
			fmt.Println(c1.Args == nil)
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	v := RPCValue{I32: 0, S:"test"}
	fmt.Println(unsafe.Sizeof(v))
	bytes := make([]byte, v.Size())
	v.MarshalTo(bytes)
	fmt.Println(len(bytes))
}
