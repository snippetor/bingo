package net

import (
	"bytes"
	"encoding/binary"
)

type ByteBuffer struct {
	buffer []byte
}

func (b *ByteBuffer) Write(bytes []byte) {
	if b.buffer == nil {
		b.buffer = make([]byte, 0)
	}
	b.buffer = append(b.buffer, bytes...)
}

func (b *ByteBuffer) Reset() {
	b.buffer = nil
}

func (b *ByteBuffer) SplitPacket() (int, []byte) {
	if b.buffer == nil || len(b.buffer) == 0 {
		return 0, nil
	}
	// 前4个字节为包长度
	var length int
	buf := bytes.NewBuffer(b.buffer[:4])
	binary.Read(buf, binary.BigEndian, &length)
	if length < 0 {
		b.Reset()
		return 0, nil
	}
	if len(b.buffer)-4 >= int(length) {
		var id int
		// 再4个字节为消息ID
		buf.Reset()
		buf.Write(b.buffer[4:8])
		binary.Read(buf, binary.BigEndian, &id)
		// 剩余为包体
		ret := b.buffer[8:4 + length]
		if int(4 + length) < len(b.buffer) {
			b.buffer = b.buffer[4 + length:]
		} else {
			b.buffer = make([]byte, 0)
		}
		return id, ret
	}
	return 0, nil
}
