package net

import (
	"encoding/binary"
	"bytes"
)

type MessagePacket struct {
	MessageId int32
	Content   []byte
}

func (p *MessagePacket) Pack() []byte {
	pk := make([]byte, 0)
	// 写入长度
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(p.Content))+4)
	pk = append(pk, buf.Bytes()...)
	// 写入id
	buf.Reset()
	binary.Write(buf, binary.BigEndian, p.MessageId)
	pk = append(pk, buf.Bytes()...)
	// 消息体
	pk = append(pk, p.Content...)
	return pk
}