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
	"net"
	"strconv"
	"sync"

	"github.com/xtaci/kcp-go"
	"github.com/snippetor/bingo/errors"
	"github.com/snippetor/bingo/log"
)

type kcpServer struct {
	sync.RWMutex
	listener net.Listener
	clients  map[uint32]Conn
}

func (s *kcpServer) listen(port int, callback MessageCallback) error {
	listener, err := kcp.Listen(":" + strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener
	s.clients = make(map[uint32]Conn, 0)
	log.I("Kcp server runnning on :%d", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.I(conn.RemoteAddr().String()+" %s", " kcp connect success")
		c := &kcpConn{conn: conn}
		c.setState(ConnStateConnected)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		callback(c, MsgIdConnConnect, nil)
		go s.handleConnection(c, callback)
	}
	return nil
}

// 处理消息流
func (s *kcpServer) handleConnection(conn Conn, callback MessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			log.E(err.Error())
			conn.setState(ConnStateClosed)
			callback(conn, MsgIdConnDisconnect, nil)
			s.Lock()
			delete(s.clients, conn.Identity())
			s.Unlock()
			break
		}
		byteBuffer = append(byteBuffer, buf[:l]...)
		for {
			id, body, remains := globalPacker.Unpack(byteBuffer)
			if body != nil {
				callback(conn, id, body)
			}
			if body == nil || len(remains) == 0 {
				break
			}
		}
	}
}

func (s *kcpServer) GetConnection(identity uint32) (Conn, bool) {
	s.RLock()
	defer s.RUnlock()
	if s.clients == nil {
		return nil, false
	} else {
		identity, ok := s.clients[identity]
		return identity, ok
	}
}

func (s *kcpServer) Close() {
	s.Lock()
	defer s.Unlock()
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

type kcpClient struct {
	sync.Mutex
	serverAddr string
	callback   MessageCallback
	conn       Conn
}

func (c *kcpClient) Reconnect() {
	c.connect(c.serverAddr, c.callback)
}

func (c *kcpClient) connect(serverAddr string, callback MessageCallback) error {
	c.serverAddr = serverAddr
	c.callback = callback
	conn, err := kcp.Dial(serverAddr)
	if err != nil {
		callback(nil, MsgIdConnDisconnect, nil)
		return err
	}
	defer conn.Close()
	c.conn = &kcpConn{conn: conn}
	c.conn.setState(ConnStateConnected)
	callback(c.conn, MsgIdConnConnect, nil)
	log.I("Kcp connect server ok :%s", serverAddr)
	c.handleConnection(c.conn, callback)
	return nil
}

// 处理消息流
func (c *kcpClient) handleConnection(conn Conn, callback MessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			log.E(err.Error())
			c.conn.setState(ConnStateClosed)
			callback(conn, MsgIdConnDisconnect, nil)
			c.conn = nil
			break
		}
		byteBuffer = append(byteBuffer, buf[:l]...)
		for {
			id, body, remains := globalPacker.Unpack(byteBuffer)
			if body != nil {
				callback(conn, id, body)
			}
			if body == nil || len(remains) == 0 {
				break
			}
		}
	}
}

func (c *kcpClient) Send(msgId MessageId, body MessageBody) error {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.State() == ConnStateConnected {
		return c.conn.Send(msgId, body)
	} else {
		return errors.ConnectionError(errors.ErrCodeConnect)
	}
	return nil
}

func (c *kcpClient) Close() {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
