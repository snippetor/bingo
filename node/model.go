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

package node

import (
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/proto"
	"github.com/snippetor/bingo/codec"
)

var (
	models = make(map[string]IModel)
)

type IModel interface {
	OnInit()
	OnReceiveServiceMessage(net.IConn, net.MessageId, interface{})
	OnReceiveHttpServiceRequest(*fasthttp.RequestCtx)
	SetNodeName(string)
	GetNodeName() string
	SetRPCServer(*rpc.Server)
	GetRPCServer() *rpc.Server
	PutRPCClient(string, *rpc.Client)
	GetRPCClient(string) (*rpc.Client, bool)
	PutService(string, net.IServer)
	GetService(string) (net.IServer, bool)
	SetDefaultMessageProtocol(codec.CodecType)
	GetDefaultMessageProtocol() *proto.MessageProtocol
	SetClientProtoVersion(string)
	GetClientProtoVersion() string
}

// rpc
// service
type Model struct {
	nodeName       string
	servers        map[string]net.IServer
	rpcClients     map[string]*rpc.Client
	rpcServer      *rpc.Server
	proto          *proto.MessageProtocol
	clientProtoVer string
}

func (m *Model) OnInit() {
}

func (m *Model) OnReceiveServiceMessage(conn net.IConn, msgId net.MessageId, msgBody interface{}) {
}

func (m *Model) OnReceiveHttpServiceRequest(ctx *fasthttp.RequestCtx) {
}

func (m *Model) CallRPCMethod(nodeName string, methodName string, args rpc.Args, callback rpc.RPCCallback) {
	if m.rpcServer != nil && m.rpcServer.HasEndName(nodeName) {
		m.rpcServer.Call(nodeName, methodName, args, callback)
	}
	if m.rpcClients != nil {
		if c, ok := m.rpcClients[nodeName]; ok && c != nil {
			c.Call(methodName, args, callback)
		}
	}
}

func (m *Model) CallRPCMethodNoReturn(nodeName string, methodName string, args rpc.Args) {
	if m.rpcServer != nil && m.rpcServer.HasEndName(nodeName) {
		m.rpcServer.CallNoReturn(nodeName, methodName, args)
	}
	if m.rpcClients != nil {
		if c, ok := m.rpcClients[nodeName]; ok && c != nil {
			c.CallNoReturn(methodName, args)
		}
	}
}

func (m *Model) RegisterRPCMethod(methodName string, method rpc.RPCMethod) {
	rpc.RegisterMethod(m.nodeName, methodName, method)
}

func (m *Model) SetNodeName(name string) {
	m.nodeName = name
}

func (m *Model) GetNodeName() string {
	return m.nodeName
}

func (m *Model) SetRPCServer(serv *rpc.Server) {
	m.rpcServer = serv
}

func (m *Model) GetRPCServer() *rpc.Server {
	return m.rpcServer
}

func (m *Model) PutRPCClient(name string, client *rpc.Client) {
	if m.rpcClients == nil {
		m.rpcClients = make(map[string]*rpc.Client)
	}
	m.rpcClients[name] = client
}

func (m *Model) GetRPCClient(name string) (*rpc.Client, bool) {
	if m.rpcClients != nil {
		if c, ok := m.rpcClients[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

func (m *Model) PutService(name string, s net.IServer) {
	if m.servers == nil {
		m.servers = make(map[string]net.IServer)
	}
	m.servers[name] = s
}

func (m *Model) GetService(name string) (net.IServer, bool) {
	if m.servers != nil {
		if c, ok := m.servers[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

func (m *Model) SetDefaultMessageProtocol(c codec.CodecType) {
	m.proto = proto.NewMessageProtocol(c)
}

func (m *Model) GetDefaultMessageProtocol() *proto.MessageProtocol {
	return m.proto
}

func (m *Model) SetClientProtoVersion(version string) {
	m.clientProtoVer = version
}

func (m *Model) GetClientProtoVersion() string {
	return m.clientProtoVer
}

func BindNodeModel(modelName string, m IModel) {
	models[modelName] = m
}

func getNodeModel(modelName string) (IModel, bool) {
	m, ok := models[modelName]
	return m, ok
}
