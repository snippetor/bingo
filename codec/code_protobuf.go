package codec

import (
	"github.com/gogo/protobuf/proto"
	"github.com/snippetor/bingo/net"
)

// Protobuf消息协议
type protoBuf struct {
}

func (p *protoBuf) Marshal(v interface{}) (net.MessageBody, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *protoBuf) Unmarshal(data net.MessageBody, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
