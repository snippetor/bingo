package protocol

import (
	"github.com/snippetor/bingo/net"
	"errors"
	"strconv"
)

type MessageProtocol byte

var (
	msgProtoType MessageProtocol
	msgProto     iProtocol
	protoVersion string
)

const (
	Json     MessageProtocol = iota
	Protobuf
)

func init() {
	SetDefaultProtocol(Protobuf)
}

// 设置默认消息协议，默认为Protobuf
func SetDefaultProtocol(proto MessageProtocol) {
	msgProtoType = proto
	switch proto {
	case Json:
		msgProto = iProtocol(&protocolJson{})
	case Protobuf:
		msgProto = iProtocol(&protocolProtoBuf{})
	}
}

// 获取默认消息协议
func GetDefaultProtocol() MessageProtocol {
	return msgProtoType
}

// 设置协议版本
func SetDefaultProtoVersion(version string) {
	protoVersion = version
}

// 获取协议版本
func GetDefaultProtoVersion() string {
	return protoVersion
}

// 根据默认协议将消息结构体转换为[]byte
// 注意：此方法要传入指针
func Marshal(v interface{}) (net.MessageBody, error) {
	return msgProto.marshal(v)
}

// 尝试使用注册的消息ID和结构体对来解析消息体，如果解析无误将返回对应的消息结构体
// 注意：此方法返回的是消息结构体的指针，而非值
func Unmarshal(msgId net.MessageId, data net.MessageBody, collection IProtoCollection) (interface{}, error) {
	if v, ok := collection.Get(msgId, protoVersion); ok {
		if err := msgProto.unmarshal(data, v); err != nil {
			return nil, err
		}
		return v, nil
	} else if v, ok := collection.GetDefault(msgId); ok {
		if err := msgProto.unmarshal(data, v); err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, errors.New("no found message struct for msgid #" + strconv.Itoa(int(msgId)))
}
