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
	"net/http"
	"strconv"
	"github.com/snippetor/bingo/comm"
	"sync"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/log/fwlogger"
)

type wsConn struct {
	conn *websocket.Conn
	absConn
}

func (c *wsConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.WriteMessage(websocket.BinaryMessage, GetDefaultMessagePacker().Pack(msgId, body))
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

func (c *wsConn) GetNetProtocol() NetProtocol {
	return WebSocket
}

type wsServer struct {
	comm.Configable
	sync.RWMutex
	upgrader *websocket.Upgrader
	callback IMessageCallback
	clients  map[utils.Identity]IConn
}

func (s *wsServer) wsHttpHandle(w http.ResponseWriter, r *http.Request) {
	if conn, err := s.upgrader.Upgrade(w, r, nil); err == nil {
		c := IConn(&wsConn{conn: conn})
		c.setState(STATE_CONNECTED)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		s.callback(c, MSGID_CONNECT_CONNECTED, nil)
		go s.handleConnection(c, s.callback)
	} else {
		fwlogger.E("-- ws build connection failed!!! --")
	}
}

func (s *wsServer) listen(port int, callback IMessageCallback) bool {
	if s.upgrader == nil {
		s.upgrader = &websocket.Upgrader{}
	}
	s.callback = callback
	s.clients = make(map[utils.Identity]IConn, 0)
	http.HandleFunc("/", s.wsHttpHandle)
	if err := http.ListenAndServe("localhost:"+strconv.Itoa(port), nil); err != nil {
		fwlogger.E(err.Error())
		return false
	}
	return true
}

func (s *wsServer) Close() {
}

func (s *wsServer) handleConnection(conn IConn, callback IMessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			fwlogger.E(err.Error())
			conn.setState(STATE_CLOSED)
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			s.Lock()
			delete(s.clients, conn.Identity())
			s.Unlock()
			break
		}
		packer := GetDefaultMessagePacker()
		id, body, _ := packer.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (s *wsServer) GetConnection(identity utils.Identity) (IConn, bool) {
	s.RLock()
	defer s.RUnlock()
	if s.clients == nil {
		return nil, false
	} else {
		identity, ok := s.clients[identity]
		return identity, ok
	}
}

type wsClient struct {
	comm.Configable
	sync.Mutex
	conn IConn
}

func (c *wsClient) connect(serverAddr string, callback IMessageCallback) bool {
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	fwlogger.I("Ws connect server ok :%s", serverAddr)
	if err != nil {
		fwlogger.E(err.Error())
		return false
	}
	c.conn = IConn(&wsConn{conn: conn})
	c.conn.setState(STATE_CONNECTED)
	callback(c.conn, MSGID_CONNECT_CONNECTED, nil)
	c.handleConnection(c.conn, callback)
	return true
}

func (c *wsClient) handleConnection(conn IConn, callback IMessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			fwlogger.E(err.Error())
			c.conn.setState(STATE_CLOSED)
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			c.conn = nil
			break
		}
		packer := GetDefaultMessagePacker()
		id, body, _ := packer.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (c *wsClient) Send(msgId MessageId, body MessageBody) bool {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.GetState() == STATE_CONNECTED {
		return c.conn.Send(msgId, body)
	} else {
		fwlogger.W("-- send tcp message failed!!! conn wrong state --")
	}
	return false
}

func (c *wsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
