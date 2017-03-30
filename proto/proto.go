package proto

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/codec"
	"github.com/snippetor/bingo/log/fwlogger"
)

type MessageProtocol struct {
	c codec.ICodec
	p IProtoCollection
}

func NewMessageProtocol(c codec.CodecType) *MessageProtocol {
	p := &MessageProtocol{}
	p.c = codec.NewCodec(c)
	p.p = IProtoCollection(&DefaultProtoCollection{})
	return p
}

func (mp *MessageProtocol) SetCodec(c codec.ICodec) {
	mp.c = c
}

func (mp *MessageProtocol) GetCodec() codec.ICodec {
	return mp.c
}

func (mp *MessageProtocol) SetProtoCollection(c IProtoCollection) {
	mp.p = c
}

func (mp *MessageProtocol) GetProtoCollection() IProtoCollection {
	return mp.p
}

// 根据默认协议将消息结构体转换为[]byte
// 注意：此方法要传入指针
func (mp *MessageProtocol) Marshal(v interface{}) (net.MessageBody, bool) {
	if b, err := mp.c.Marshal(v); err == nil {
		return b, true
	}
	fwlogger.E("-- proto Marshal message failed! --")
	return nil, false
}

// 尝试使用注册的消息ID和结构体对来解析消息体，如果解析无误将返回对应的消息结构体
// 注意：此方法返回的是消息结构体的指针，而非值
// @protoVersion 客服端传过来的协议版本
func (mp *MessageProtocol) UnmarshalTo(msgId net.MessageId, data net.MessageBody, clientProtoVersion string) (interface{}, bool) {
	if v, ok := mp.p.Get(msgId, clientProtoVersion); ok {
		if err := mp.c.Unmarshal(data, v); err != nil {
			fwlogger.E("-- proto UnmarshalTo message failed! --")
			return nil, false
		}
		return v, true
	} else if v, ok := mp.p.GetDefault(msgId); ok {
		if err := mp.c.Unmarshal(data, v); err != nil {
			fwlogger.E("-- proto UnmarshalTo message failed! --")
			return nil, false
		}
		return v, true
	}
	fwlogger.E("-- proto UnmarshalTo message failed! no found message struct for msgid #%d --", int(msgId))
	return nil, false
}

func (mp *MessageProtocol) UnmarshalToDefault(msgId net.MessageId, data net.MessageBody) (interface{}, bool) {
	if v, ok := mp.p.GetDefault(msgId); ok {
		if err := mp.c.Unmarshal(data, v); err != nil {
			fwlogger.E("-- proto UnmarshalTo message failed! --")
			return nil, false
		}
		return v, true
	}
	fwlogger.E("-- proto UnmarshalTo message failed! no found message struct for msgid #%d --", int(msgId))
	return nil, false
}

// 直接解析成结构体
func (mp *MessageProtocol) Unmarshal(data net.MessageBody, v interface{}) bool {
	if err := mp.c.Unmarshal(data, v); err != nil {
		fwlogger.E("-- proto Unmarshal message failed! --")
		return false
	}
	return true
}
