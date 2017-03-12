package protocol

import (
	"github.com/golang/protobuf/proto"
	"github.com/snippetor/bingo/net"
)

// Protobuf消息协议
type protocolProtoBuf struct {
}

func (p *protocolProtoBuf) marshal(v interface{}) (net.MessageBody, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *protocolProtoBuf) unmarshal(data net.MessageBody, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
