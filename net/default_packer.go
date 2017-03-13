package net

import (
	"encoding/binary"
	"bytes"
	"math"
	"github.com/snippetor/bingo"
)

// 消息包
type DefaultMessagePacker struct {
}

// TODO 性能优化，目前 3000000	       541 ns/op
func (p *DefaultMessagePacker) Pack(messageId MessageId, body MessageBody) []byte {
	bodyLen := len(body)
	if body != nil && bodyLen > math.MaxInt32 {
		bingo.E("-- MessagePacker - body length is too large! %d --", len(body))
		return nil
	}
	pk := make([]byte, 0)
	// 写入长度
	buf := bytes.NewBuffer([]byte{})
	if body == nil {
		binary.Write(buf, binary.BigEndian, 4)
	} else {
		binary.Write(buf, binary.BigEndian, int32(bodyLen)+4)
	}
	pk = append(pk, buf.Bytes()...)
	// 写入id
	buf.Reset()
	binary.Write(buf, binary.BigEndian, messageId)
	pk = append(pk, buf.Bytes()...)
	// 消息体
	if body != nil && bodyLen > 0 {
		pk = append(pk, body...)
	}
	return pk
}

func (p *DefaultMessagePacker) Unpack(buffer []byte) (MessageId, MessageBody, []byte) {
	if buffer == nil || len(buffer) == 0 {
		return -1, nil, buffer
	}
	// 前4个字节为包长度
	var length int32
	buf := bytes.NewBuffer(buffer[:4])
	binary.Read(buf, binary.BigEndian, &length)
	if length < 0 {
		return -1, nil, buffer
	}
	if len(buffer)-4 >= int(length) {
		var id MessageId
		// 再4个字节为消息ID
		buf.Reset()
		buf.Write(buffer[4:8])
		binary.Read(buf, binary.BigEndian, &id)
		// 剩余为包体
		ret := buffer[8:4 + length]
		if int(4 + length) < len(buffer) {
			buffer = buffer[4 + length:]
		} else {
			buffer = make([]byte, 0)
		}
		return id, ret, buffer
	}
	return -1, nil, buffer
}
