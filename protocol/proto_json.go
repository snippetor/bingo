package protocol

import (
	"encoding/json"
	"github.com/snippetor/bingo/net"
)

// JSON消息协议
type protocolJson struct {
}

func (j *protocolJson) marshal(v interface{}) (net.MessageBody, error) {
	return json.Marshal(v)
}

func (j *protocolJson) unmarshal(data net.MessageBody, v interface{}) error {
	return json.Unmarshal(data, v)
}
