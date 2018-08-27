package codec

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
	p := json{}
	bs := p.Marshal(&JsonObj{Id: 1, Name: "carl"})

	var b JsonObj
	p.Unmarshal(bs, &b)
	fmt.Println(b)
}

func TestProtocolProtobuf(t *testing.T) {
	p := protoBuf{}
	persion := &Person{}
	persion.Id = 1
	persion.Name = "carl"
	bs := p.Marshal(persion)

	persion2 := Person{}
	p.Unmarshal(bs, &persion2)
	fmt.Println(persion2.Name)
}

func TestGoroutines(t *testing.T) {
	p := json{}
	persion := &Person{}
	persion.Id = 1
	persion.Name = "carl"
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println(p.Marshal(persion))
		}()
	}
}
