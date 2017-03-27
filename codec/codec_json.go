package codec

import (
	json1 "encoding/json"
	"github.com/snippetor/bingo/net"
)

// JSON消息协议
type json struct {
}

func (j *json) Marshal(v interface{}) (net.MessageBody, error) {
	return json1.Marshal(v)
}

func (j *json) Unmarshal(data net.MessageBody, v interface{}) error {
	return json1.Unmarshal(data, v)
}
