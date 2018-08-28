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
	"errors"
	"github.com/snippetor/bingo/utils"
)

// 网络连接接口
type Conn interface {
	Send(msgId MessageId, body MessageBody) error
	Close()
	Address() string
	read(*[]byte) (int, error)
	NetProtocol() Protocol
	Identity() uint32
	State() ConnState
	setState(ConnState)
}

type ConnState int

const (
	ConnStateClosed     ConnState = iota
	ConnStateConnecting
	ConnStateConnected
)

// Conn ID生成器
var (
	identifier *utils.Identifier
)

func init() {
	identifier = utils.NewIdentifier(1)
}

type baseConn struct {
	identity uint32
	state    ConnState
}

func (c *baseConn) Send(msgId MessageId, body MessageBody) error {
	return errors.New("not implements: send")
}

func (c *baseConn) Close() {
}

func (c *baseConn) Address() string {
	return "0.0.0.0"
}

func (c *baseConn) read(*[]byte) (int, error) {
	return -1, errors.New("not implements: read")
}

func (c *baseConn) NetProtocol() Protocol {
	return -1
}

func (c *baseConn) Identity() uint32 {
	if !identifier.IsValidIdentity(c.identity) {
		c.identity = identifier.GenIdentity()
	}
	return c.identity
}

func (c *baseConn) State() ConnState {
	return c.state
}

func (c *baseConn) setState(state ConnState) {
	c.state = state
}
