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
	"io/ioutil"
	"encoding/json"
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"strings"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/snippetor/bingo/log/fwlogger"
)

type Service struct {
	Name  string
	Type  string
	Port  int
	Codec string
}

type Node struct {
	Name      string
	ModelName string `json:"model"`
	Domain    int
	Services  []*Service `json:"service"`
	RpcPort   int        `json:"rpc-port"`
	RpcTo     []string   `json:"rpc-to"`
}

type Config struct {
	Domains []string   `json:"domains"`
	Nodes   []*Node    `json:"node"`
}

var (
	config *Config
	nodes  = make(map[string]*IModel, 0)
)

func init() {
	config = &Config{}
}

func Parse(configPath string) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
	if err = json.Unmarshal(content, config); err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
	//TODO check
}

func findNode(name string) *Node {
	for _, n := range config.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func run_node(n *Node) {
	// create node
	m, ok := getNodeModel(n.ModelName)
	if !ok {
		fwlogger.E("-- not found model by name %s --", n.ModelName)
		return
	}
	m.setNodeName(n.Name)
	m.init()
	// rpc
	if n.RpcPort > 0 {
		s := &rpc.Server{}
		s.Listen(n.Name, n.RpcPort)
		m.setRPCServer(s)
	}
	for _, rpcServerNode := range n.RpcTo {
		c := &rpc.Client{}
		serverNode := findNode(rpcServerNode)
		if serverNode != nil {
			c.Connect(n.Name, config.Domains[serverNode.Domain]+":"+strconv.Itoa(serverNode.RpcPort))
			m.putRPCClient(serverNode.Name, c)
		}
	}
	// service
	for _, s := range n.Services {
		switch strings.ToLower(s.Type) {
		case "tcp":
			serv := net.GoListen(net.Tcp, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				switch msgId {
				case net.MSGID_CONNECT_CONNECTED:
					m.OnServiceClientConnected(s.Name, conn)
				case net.MSGID_CONNECT_DISCONNECT:
					m.OnServiceClientDisconnected(s.Name, conn)
				default:
					m.OnReceiveServiceMessage(conn, msgId, body)
				}
			})
			m.putService(s.Name, serv)
		case "ws":
			serv := net.GoListen(net.WebSocket, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				switch msgId {
				case net.MSGID_CONNECT_CONNECTED:
					m.OnServiceClientConnected(s.Name, conn)
				case net.MSGID_CONNECT_DISCONNECT:
					m.OnServiceClientDisconnected(s.Name, conn)
				default:
					m.OnReceiveServiceMessage(conn, msgId, body)
				}
			})
			m.putService(s.Name, serv)
		case "http":
			go func() {
				if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(ctx *fasthttp.RequestCtx) {
					m.OnReceiveHttpServiceRequest(ctx)
				}); err != nil {
					fwlogger.E("-- startup http service failed! %s --", err)
				}
			}()
		}
	}

	nodes[n.Name] = &m
}

func Run(nodeName string) {
	n := findNode(nodeName)
	if n == nil {
		fwlogger.E("-- run node failed! not found node by name %s --", nodeName)
		return
	}
	run_node(n)
}

func Stop(nodeName string) {
	n := findNode(nodeName)
	if n == nil {
		fwlogger.E("-- stop node failed! not found node by name %s --", nodeName)
		return
	}
	if m, ok := nodes[nodeName]; ok {
		(*m).destroy()
	}
}

func RunAll() {
	for _, n := range config.Nodes {
		run_node(n)
	}
}

func StopAll() {
	for _, m := range nodes {
		if m != nil {
			(*m).destroy()
		}
	}
}
