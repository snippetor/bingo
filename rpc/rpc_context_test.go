package rpc

import (
	"testing"
	"math"
	"fmt"
)

func TestArgs(t *testing.T) {
	a := make(Args)
	a.Put("test1", "test1")
	a.PutInt("test2", math.MaxInt32)
	a.PutInt64("test3", math.MaxInt64)
	a.PutFloat32("test4", math.MaxFloat32)
	a.PutBool("test5", true)

	if v, ok := a.Get("test1"); !ok || v != "test1" {
		t.Fail()
		fmt.Println(a)
	}
	if a.MustGet("test1", "default") != "test1" {
		t.Fail()
		fmt.Println(a)
	}

	if v, ok := a.GetInt("test2"); !ok || v != math.MaxInt32 {
		t.Fail()
		fmt.Println(a)
	}
	if a.MustGetInt("test2", 1) != math.MaxInt32 {
		t.Fail()
		fmt.Println(a)
	}

	if v, ok := a.GetInt64("test3"); !ok || v != math.MaxInt64 {
		t.Fail()
		fmt.Println(a)
	}
	if a.MustGetInt64("test3", 1) != math.MaxInt64 {
		t.Fail()
		fmt.Println(a)
	}

	if v, ok := a.GetFloat32("test4"); !ok || v != math.MaxFloat32 {
		t.Fail()
		fmt.Println(a)
	}
	fmt.Println(a.MustGetFloat32("test4", 1))
	if a.MustGetFloat32("test4", 1) != math.MaxFloat32 {
		t.Fail()
		fmt.Println(a)
	}

	if v, ok := a.GetBool("test5"); !ok || v != true {
		t.Fail()
		fmt.Println(a)
	}
	if a.MustGetBool("test5", false) != true {
		t.Fail()
		fmt.Println(a)
	}

}
