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
	"strings"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/log"
)

var (
	models = make(map[string]IModel)
)

type IModel interface {
	OnInit()
	OnDestroy()
	OnServiceClientConnected(string, net.IConn)
	OnServiceClientDisconnected(string, net.IConn)
	OnReceiveServiceMessage(net.IConn, net.MessageId, net.MessageBody)
	OnReceiveHttpServiceRequest(*fasthttp.RequestCtx)

	init()
	setNodeName(string)
	setRPCServer(*rpc.Server)
	putRPCClient(string, *rpc.Client)
	putService(string, net.IServer)
	setConfig(*utils.ValueMap)
	setLoggers(map[string]*log.Logger)
	destroy()
}

// rpc
type RPCModule struct {
	nodeName   string
	rpcClients []*rpc.Client
	rpcServer  *rpc.Server
	rpcRouter  *rpc.Router
}

func (m *RPCModule) RegisterController(c rpc.IController) {
	m.rpcRouter.RegisterController(m.nodeName, c)
}

func (m *RPCModule) GetEndStub(nodeName string) (rpc.IEndStub, bool) {
	if m.rpcServer != nil {
		if c, ok := m.rpcServer.GetClient(nodeName); ok {
			return rpc.IEndStub(c), true
		}
	}
	if m.rpcClients != nil && len(m.rpcClients) > 0 {
		for _, c := range m.rpcClients {
			if c.EndName == nodeName {
				return rpc.IEndStub(c), true
			}
		}
	}
	return nil, false
}

func (m *RPCModule) GetEndStubsWithPrefix(nodeNamePrefix string) []rpc.IEndStub {
	stubs := make([]rpc.IEndStub, 0)
	if m.rpcServer != nil {
		for _, c := range m.rpcServer.GetClients() {
			if strings.HasPrefix(c.EndName, nodeNamePrefix) {
				stubs = append(stubs, rpc.IEndStub(c))
			}
		}
	}
	if m.rpcClients != nil {
		for _, c := range m.rpcClients {
			if strings.HasPrefix(c.EndName, nodeNamePrefix) {
				stubs = append(stubs, rpc.IEndStub(c))
			}
		}
	}
	return stubs
}

func (m *RPCModule) GetEndStubsWithModelName(nodeModelName string) []rpc.IEndStub {
	stubs := make([]rpc.IEndStub, 0)
	if m.rpcServer != nil {
		for _, c := range m.rpcServer.GetClients() {
			if nodeModelName == c.EndModelName {
				stubs = append(stubs, rpc.IEndStub(c))
			}
		}
	}
	if m.rpcClients != nil {
		for _, c := range m.rpcClients {
			if nodeModelName == c.EndModelName {
				stubs = append(stubs, rpc.IEndStub(c))
			}
		}
	}
	return stubs
}

// 通过名称前缀和平衡因子获取stub
// 如果有多个stub，则通过取模算法（balancingSeed % (stubs size)）来决定使用哪个
func (m *RPCModule) GetEndStubWithPrefixAndBalancingSeed(nodeNamePrefix string, balancingSeed int64) (rpc.IEndStub, bool) {
	stubs := m.GetEndStubsWithPrefix(nodeNamePrefix)
	len := len(stubs)
	if len == 0 {
		return nil, false
	} else if len == 1 {
		return stubs[0], true
	} else {
		return stubs[balancingSeed%int64(len)], true
	}
}

// 通过节点模型名称和平衡因子获取stub
// 如果有多个stub，则通过取模算法（balancingSeed % (stubs size)）来决定使用哪个
func (m *RPCModule) GetEndStubWithModelNameAndBalancingSeed(nodeModelName string, balancingSeed int64) (rpc.IEndStub, bool) {
	stubs := m.GetEndStubsWithModelName(nodeModelName)
	len := len(stubs)
	if len == 0 {
		return nil, false
	} else if len == 1 {
		return stubs[0], true
	} else {
		return stubs[balancingSeed%int64(len)], true
	}
}

func (m *RPCModule) Close() {
	if m.rpcServer != nil {
		m.rpcServer.Close()
	}
	if m.rpcClients != nil {
		for _, v := range m.rpcClients {
			if v != nil {
				v.Close()
			}
		}
	}
}

// service
type ServiceModule struct {
	servers map[string]net.IServer
}

func (m *ServiceModule) GetService(name string) (net.IServer, bool) {
	if m.servers != nil {
		if c, ok := m.servers[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

func (m *ServiceModule) Close() {
	if m.servers != nil {
		for _, v := range m.servers {
			if v != nil {
				v.Close()
			}
		}
	}
}

// log module
type LogModule struct {
	loggers map[string]*log.Logger
}

func (m *LogModule) GetLogger(name string) *log.Logger {
	if m.loggers != nil {
		if logger, ok := m.loggers[name]; ok {
			return logger
		}
	}
	return nil
}

// model
type Model struct {
	nodeName string
	LOG      *LogModule
	RPC      *RPCModule
	Service  *ServiceModule
	Config   *utils.ValueMap
}

func (m *Model) init() {
	// init Log module
	m.LOG = &LogModule{}
	// init RPC module
	m.RPC = &RPCModule{nodeName: m.nodeName, rpcRouter: &rpc.Router{}}
	// init Service module
	m.Service = &ServiceModule{}
}

func (m *Model) setNodeName(name string) {
	m.nodeName = name
}

func (m *Model) putService(name string, s net.IServer) {
	if m.Service.servers == nil {
		m.Service.servers = make(map[string]net.IServer)
	}
	m.Service.servers[name] = s
}

func (m *Model) setRPCServer(serv *rpc.Server) {
	serv.SetRouter(m.RPC.rpcRouter)
	m.RPC.rpcServer = serv
}

func (m *Model) putRPCClient(name string, client *rpc.Client) {
	client.SetRouter(m.RPC.rpcRouter)
	m.RPC.rpcClients = append(m.RPC.rpcClients, client)
}

func (m *Model) setConfig(config *utils.ValueMap) {
	m.Config = config
}

func (m *Model) setLoggers(loggers map[string]*log.Logger) {
	m.LOG.loggers = loggers
}

func (m *Model) destroy() {
	// close RPC
	m.OnDestroy()
}

func (m *Model) OnInit() {
}

func (m *Model) OnDestroy() {
}

func (m *Model) OnServiceClientConnected(serviceName string, conn net.IConn) {
}

func (m *Model) OnServiceClientDisconnected(serviceName string, conn net.IConn) {
}

func (m *Model) OnReceiveServiceMessage(conn net.IConn, msgId net.MessageId, msgBody net.MessageBody) {
}

func (m *Model) OnReceiveHttpServiceRequest(ctx *fasthttp.RequestCtx) {
}

func (m *Model) GetNodeName() string {
	return m.nodeName
}

func BindNodeModel(modelName string, m IModel) {
	models[modelName] = m
}

func getNodeModel(modelName string) (IModel, bool) {
	m, ok := models[modelName]
	return m, ok
}
