package proto

import "github.com/snippetor/bingo/net"

// 协议集合，用于存储消息ID和结构体的对应关系集合
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
