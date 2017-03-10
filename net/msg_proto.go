package net

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

type protocolJson struct {
}

func (j *protocolJson) marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *protocolJson) unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type protocolProtoBuf struct {
}

func (p *protocolProtoBuf) marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *protocolProtoBuf) unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
