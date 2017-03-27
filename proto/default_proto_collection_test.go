package proto

import (
	"testing"
	"fmt"
	"reflect"
)

type TestObj struct {
	Id   int
	Name string
}

func TestProtoCollection(t *testing.T) {
	c := DefaultProtoCollection{}
	c.PutDefault(1, TestObj{})

	c1, ok := c.GetDefault(1)
	if ok {
		o := c1.(*TestObj)
		o.Id = 1
	}

	c2, ok := c.GetDefault(1)
	if ok {
		o := c2.(*TestObj)
		fmt.Println(o.Id)
	}

	c.PutDefault(2, &TestObj{})
	c3, ok := c.GetDefault(2)
	if ok {
		o := c3.(*TestObj)
		fmt.Println(o.Id)
	}

	c.Put(2, &TestObj1{}, "1.1")
	c4, ok := c.Get(2, "1.1")
	if ok {
		o := c4.(*TestObj1)
		fmt.Println(o.Id)
	}

	if c.Size() != 3 {
		t.Fail()
	}
}

type TestObj1 struct {
	Id      int
	Name    string
	Avatar  string
	Phone   string
	Id1     int
	Name1   string
	Avatar1 string
	Phone1  string
	Id2     int
	Name2   string
	Avatar2 string
	Phone2  string
	Id3     int
	Name3   string
	Avatar3 string
	Phone3  string
	Id4     int
	Name4   string
	Avatar4 string
	Phone4  string
}

func BenchmarkReflectNew(b *testing.B) {
	b.StopTimer()
	t := reflect.TypeOf(TestObj1{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		reflect.New(t).Interface()
	}
}
