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

func (p *DefaultMessagePacker) Pack(messageId MessageId, content MessageBody) []byte {
	if content != nil && len(content) > math.MaxInt32 {
		bingo.E("-- MessagePacker - body length is too large! %d --", len(content))
		return nil
	}
	pk := make([]byte, 0)
	// 写入长度
	buf := bytes.NewBuffer([]byte{})
	if content == nil {
		binary.Write(buf, binary.BigEndian, 4)
	} else {
		binary.Write(buf, binary.BigEndian, int32(len(content))+4)
	}
	pk = append(pk, buf.Bytes()...)
	// 写入id
	buf.Reset()
	binary.Write(buf, binary.BigEndian, messageId)
	pk = append(pk, buf.Bytes()...)
	// 消息体
	if content != nil && len(content) > 0 {
		pk = append(pk, content...)
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
