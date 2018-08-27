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
	"github.com/gorilla/websocket"
	"github.com/snippetor/bingo/log/fwlogger"
)

type wsConn struct {
	baseConn
	conn *websocket.Conn
}

func (c *wsConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.WriteMessage(websocket.BinaryMessage, globalPacker.Pack(msgId, body))
		return true
	} else {
		fwlogger.W("-- send message failed!!! --")
		return false
	}
}

func (c *wsConn) read(buf *[]byte) (int, error) {
	if c.conn != nil {
		t, msg, err := c.conn.ReadMessage()
		if err == nil && t == websocket.BinaryMessage {
			*buf = msg
			return len(msg), nil
		}
	}
	return -1, nil
}

func (c *wsConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *wsConn) Address() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return "0:0:0:0"
}

func (c *wsConn) NetProtocol() Protocol {
	return WebSocket
}
