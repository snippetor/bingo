package proto

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/codec"
	"github.com/snippetor/bingo"
)

var (
	defaultCodec           codec.ICodec
	defaultProtoCollection IProtoCollection
)

func init() {
	defaultCodec = codec.NewCodec(codec.Protobuf)
	defaultProtoCollection = IProtoCollection(&DefaultProtoCollection{})
}

func SetDefaultCodec(c codec.ICodec) {
	defaultCodec = c
}

func SetDefaultProtoCollection(c IProtoCollection) {
	defaultProtoCollection = c
}

// 根据默认协议将消息结构体转换为[]byte
// 注意：此方法要传入指针
func Marshal(v interface{}) (net.MessageBody, bool) {
	if b, err := defaultCodec.Marshal(v); err == nil {
		return b, true
	}
	bingo.E("-- proto Marshal message failed! --")
	return nil, false
}

// 尝试使用注册的消息ID和结构体对来解析消息体，如果解析无误将返回对应的消息结构体
// 注意：此方法返回的是消息结构体的指针，而非值
// @protoVersion 客服端传过来的协议版本
func UnmarshalTo(msgId net.MessageId, data net.MessageBody, clientProtoVersion string) (interface{}, bool) {
	if v, ok := defaultProtoCollection.Get(msgId, clientProtoVersion); ok {
		if err := defaultCodec.Unmarshal(data, v); err != nil {
			bingo.E("-- proto UnmarshalTo message failed! --")
			return nil, false
		}
		return v, true
	} else if v, ok := defaultProtoCollection.GetDefault(msgId); ok {
		if err := defaultCodec.Unmarshal(data, v); err != nil {
			bingo.E("-- proto UnmarshalTo message failed! --")
			return nil, false
		}
		return v, true
	}
	bingo.E("-- proto UnmarshalTo message failed! no found message struct for msgid #%d --", int(msgId))
	return nil, false
}

// 直接解析成结构体
func Unmarshal(data net.MessageBody, v interface{}) bool {
	if err := defaultCodec.Unmarshal(data, v); err != nil {
		bingo.E("-- proto Unmarshal message failed! --")
		return false
	}
	return true
}
