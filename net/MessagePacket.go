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

type MessageSplitter struct {
	buffer []byte
}

func (s *MessageSplitter) Write(bytes []byte) {
	if s.buffer == nil {
		s.buffer = make([]byte, 0)
	}
	s.buffer = append(s.buffer, bytes...)
}

func (s *MessageSplitter) Reset() {
	s.buffer = nil
}

func (s *MessageSplitter) SplitPacket() (int, []byte) {
	if s.buffer == nil || len(s.buffer) == 0 {
		return 0, nil
	}
	// 前4个字节为包长度
	var length int
	buf := bytes.NewBuffer(s.buffer[:4])
	binary.Read(buf, binary.BigEndian, &length)
	if length < 0 {
		s.Reset()
		return 0, nil
	}
	if len(s.buffer)-4 >= int(length) {
		var id int
		// 再4个字节为消息ID
		buf.Reset()
		buf.Write(s.buffer[4:8])
		binary.Read(buf, binary.BigEndian, &id)
		// 剩余为包体
		ret := s.buffer[8:4 + length]
		if int(4 + length) < len(s.buffer) {
			s.buffer = s.buffer[4 + length:]
		} else {
			s.buffer = make([]byte, 0)
		}
		return id, ret
	}
	return 0, nil
}
