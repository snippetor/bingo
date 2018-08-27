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

package rpc

import (
	"time"
	"github.com/smallnest/rpcx/server"
	"strconv"
	"github.com/smallnest/rpcx/client"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/serverplugin"
	"context"
	"github.com/smallnest/rpcx/protocol"
	"github.com/gogo/protobuf/proto"
	"errors"
	"github.com/snippetor/bingo/codec"
	"github.com/snippetor/bingo/net"
)

type Server struct {
	name    string
	appName string
	serv    *server.Server
}

func (s *Server) Listen(name, appName string, port int, etcdAddrs []string, onInit func(server *Server)) {
	s.name = name
	s.appName = appName
	s.serv = server.NewServer()
	go func() {
		r := &serverplugin.EtcdRegisterPlugin{
			ServiceAddress: "kcp@:" + strconv.Itoa(port),
			EtcdServers:    etcdAddrs,
			BasePath:       "bingo",
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Minute,
		}
		err := r.Start()
		if err != nil {
			panic(err)
		}
		s.serv.Plugins.Add(r)
		if onInit != nil {
			onInit(s)
		}
		if err := s.serv.Serve("kcp", ":"+strconv.Itoa(port)); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) RegisterFunction(name string, fn interface{}) {
	s.serv.RegisterFunctionName(s.appName, name, fn, "")
}

func (s *Server) Close() {
	s.serv.Close()
}

type Client struct {
	name          string
	appName       string
	addr          string
	ServerAppName string
	client        client.XClient
}

func (c *Client) Connect(name, appName, serverAppName string, etcdAddrs []string) {
	c.name = name
	c.appName = appName
	c.ServerAppName = serverAppName
	d := client.NewEtcdDiscovery("bingo", serverAppName, etcdAddrs, nil)
	var option = client.DefaultOption
	option.SerializeType = protocol.SerializeNone
	c.client = client.NewXClient(serverAppName, client.Failtry, client.RoundRobin, d, option)
}

func (c *Client) Call(method string, args interface{}, reply interface{}) error {
	if msg, ok := args.(proto.Message); ok {
		if body, err := codec.ProtobufCodec.Marshal(msg); err == nil {
			var ret []byte
			if err := c.client.Call(context.Background(), "RPC:"+method, body, &ret); err == nil {
				if ret == nil {
					return nil
				}
				return codec.ProtobufCodec.Unmarshal(net.MessageBody(ret), reply)
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("wrong type of args, should be *proto.Message")
	}
}

func (c *Client) Close() {
	c.client.Close()
}
