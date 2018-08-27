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

	"github.com/snippetor/bingo/log/fwlogger"
)

type tcpServer struct {
	sync.RWMutex
	listener *net.TCPListener
	clients  map[uint32]Conn
}

func (s *tcpServer) listen(port int, callback MessageCallback) bool {
	addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fwlogger.E(err.Error())
		return false
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fwlogger.E(err.Error())
		return false
	}
	defer listener.Close()
	s.listener = listener
	s.clients = make(map[uint32]Conn, 0)
	fwlogger.I("Tcp server runnning on :%d", port)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		fwlogger.I(conn.RemoteAddr().String()+" %s", " tcp connect success")
		c := Conn(&tcpConn{conn: conn})
		c.setState(ConnStateConnected)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		callback(c, MsgIdConnConnect, nil)
		go s.handleConnection(c, callback)
	}
	return true
}

// 处理消息流
func (s *tcpServer) handleConnection(conn Conn, callback MessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			fwlogger.E(err.Error())
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

func (s *tcpServer) GetConnection(identity uint32) (Conn, bool) {
	s.RLock()
	defer s.RUnlock()
	if s.clients == nil {
		return nil, false
	} else {
		identity, ok := s.clients[identity]
		return identity, ok
	}
}

func (s *tcpServer) Close() {
	s.Lock()
	defer s.Unlock()
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

type tcpClient struct {
	sync.Mutex
	serverAddr string
	callback   MessageCallback
	conn       Conn
}

func (c *tcpClient) Reconnect() {
	c.connect(c.serverAddr, c.callback)
}

func (c *tcpClient) connect(serverAddr string, callback MessageCallback) bool {
	c.serverAddr = serverAddr
	c.callback = callback
	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		fwlogger.E(err.Error())
		callback(nil, MsgIdConnDisconnect, nil)
		return false
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fwlogger.E(err.Error())
		callback(nil, MsgIdConnDisconnect, nil)
		return false
	}
	defer conn.Close()
	c.conn = Conn(&tcpConn{conn: conn})
	c.conn.setState(ConnStateConnected)
	callback(c.conn, MsgIdConnConnect, nil)
	fwlogger.I("Tcp connect server ok :%s", serverAddr)
	c.handleConnection(c.conn, callback)
	return true
}

// 处理消息流
func (c *tcpClient) handleConnection(conn Conn, callback MessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			fwlogger.E(err.Error())
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

func (c *tcpClient) Send(msgId MessageId, body MessageBody) bool {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.State() == ConnStateConnected {
		return c.conn.Send(msgId, body)
	} else {
		fwlogger.W("-- send tcp message failed!!! conn wrong state --")
	}
	return false
}

func (c *tcpClient) Close() {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
