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
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/snippetor/bingo/log"
	"github.com/snippetor/bingo/errors"
)

type wsServer struct {
	sync.RWMutex
	upgrader *websocket.Upgrader
	callback MessageCallback
	clients  map[uint32]Conn
}

func (s *wsServer) wsHttpHandle(w http.ResponseWriter, r *http.Request) {
	if conn, err := s.upgrader.Upgrade(w, r, nil); err == nil {
		c := &wsConn{conn: conn}
		c.setState(ConnStateConnected)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		s.callback(c, MsgIdConnConnect, nil)
		go s.handleConnection(c, s.callback)
	} else {
		log.E("ws build connection failed!")
	}
}

func (s *wsServer) listen(port int, callback MessageCallback) error {
	if s.upgrader == nil {
		s.upgrader = &websocket.Upgrader{}
	}
	s.callback = callback
	s.clients = make(map[uint32]Conn, 0)
	http.HandleFunc("/", s.wsHttpHandle)
	if err := http.ListenAndServe("localhost:"+strconv.Itoa(port), nil); err != nil {
		return err
	}
	return nil
}

func (s *wsServer) Close() {
}

func (s *wsServer) handleConnection(conn Conn, callback MessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			log.E(err.Error())
			conn.setState(ConnStateClosed)
			callback(conn, MsgIdConnDisconnect, nil)
			s.Lock()
			delete(s.clients, conn.Identity())
			s.Unlock()
			break
		}
		id, body, _ := globalPacker.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (s *wsServer) GetConnection(identity uint32) (Conn, bool) {
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
	sync.Mutex
	serverAddr string
	callback   MessageCallback
	conn       Conn
}

func (c *wsClient) Reconnect() {
	c.connect(c.serverAddr, c.callback)
}

func (c *wsClient) connect(serverAddr string, callback MessageCallback) error {
	c.serverAddr = serverAddr
	c.callback = callback
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		callback(nil, MsgIdConnDisconnect, nil)
		return err
	}
	c.conn = Conn(&wsConn{conn: conn})
	c.conn.setState(ConnStateConnected)
	callback(c.conn, MsgIdConnConnect, nil)
	log.I("Ws connect server ok :%s", serverAddr)
	c.handleConnection(c.conn, callback)
	return nil
}

func (c *wsClient) handleConnection(conn Conn, callback MessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			log.E(err.Error())
			c.conn.setState(ConnStateClosed)
			callback(conn, MsgIdConnDisconnect, nil)
			c.conn = nil
			break
		}
		id, body, _ := globalPacker.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (c *wsClient) Send(msgId MessageId, body MessageBody) error {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.State() == ConnStateConnected {
		return c.conn.Send(msgId, body)
	} else {
		return errors.ConnectionError(errors.ErrCodeConnect)
	}
	return nil
}

func (c *wsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
