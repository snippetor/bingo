package codec

import "github.com/snippetor/bingo/net"

// 数据传输协议，即包体格式
type ICodec interface {
	Marshal(interface{}) (net.MessageBody, error)
	Unmarshal(net.MessageBody, interface{}) error
}

type IProtoCollection interface {
	PutDefault(id net.MessageId, v interface{})
	Put(id net.MessageId, v interface{}, protoVersion string)
	GetDefault(id net.MessageId) (interface{}, bool)
	Get(id net.MessageId, protoVersion string) (interface{}, bool)
	RemoveDefault(id net.MessageId)
	Remove(id net.MessageId, protoVersion string)
	Clear()
	Size() int
	//Dump()
}
