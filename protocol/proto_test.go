package protocol

import (
	"testing"
	"fmt"
)

type JsonObj struct {
	Id       int
	Name     string
	Sex      int
	Address  string
	Phone    string
	Avatar   string
	Id1      int
	Name1    string
	Sex1     int
	Address1 string
	Phone1   string
	Avatar1  string
	Id2      int
	Name2    string
	Sex2     int
	Address2 string
	Phone2   string
	Avatar2  string
}

func TestProtocolJson(t *testing.T) {
	p := protocolJson{}
	bs, err := p.marshal(&JsonObj{Id: 1, Name: "carl"})
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
	p := protocolProtoBuf{}
	persion := &Person{}
	persion.Id = 1
	persion.Name = "carl"
	bs, err := p.marshal(persion)
	if err != nil {
		t.Fail()
	}

	persion2 := Person{}
	err = p.unmarshal(bs, &persion2)
	if err != nil {
		t.Fail()
	}
	fmt.Println(persion2.Name)
}

func TestProto(t *testing.T) {
	c := DefaultProtoCollection{}
	c.PutDefault(1, JsonObj{})
	c.PutDefault(2, Person{})

	SetDefaultProtocol(Json)
	var bytes []byte
	var err error
	if bytes, err = Marshal(&JsonObj{Id: 1, Name: "carl"}); err != nil {
		t.Fail()
	}
	var o interface{}
	if o, err = UnmarshalTo(1, bytes, IProtoCollection(&c)); err != nil {
		t.Fail()
	}

	fmt.Println(o.(*JsonObj))

	SetDefaultProtocol(Protobuf)
	if bytes, err = Marshal(&Person{Id: 1, Name: "carl"}); err != nil {
		t.Fail()
	}
	if o, err = UnmarshalTo(2, bytes, IProtoCollection(&c)); err != nil {
		t.Fail()
	}

	fmt.Println(o.(*Person))
}
