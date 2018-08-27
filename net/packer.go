// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"encoding/binary"
	"bytes"
	"math"
	"github.com/snippetor/bingo/log/fwlogger"
)

// 消息封装器接口
type MessagePacker interface {
	// 封包，传入消息ID和包体，返回字节集
	Pack(MessageId, MessageBody) []byte
	// 解包，传入符合包结构的字节集，返回消息ID，包体，剩余内容
	Unpack([]byte) (MessageId, MessageBody, []byte)
}

// 消息包
type messagePacker struct {
}

// 3000000	       541 ns/op
func (p *messagePacker) Pack(messageId MessageId, body MessageBody) []byte {
	bodyLen := len(body)
	if body != nil && bodyLen > math.MaxInt32 {
		fwlogger.E("-- MessagePacker - body length is too large! %d --", len(body))
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

// 10000000	       157 ns/op
func (p *messagePacker) Unpack(buffer []byte) (MessageId, MessageBody, []byte) {
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
		ret := buffer[8 : 4+length]
		if int(4+length) < len(buffer) {
			buffer = buffer[4+length:]
		} else {
			buffer = make([]byte, 0)
		}
		return id, ret, buffer
	}
	return -1, nil, buffer
}
