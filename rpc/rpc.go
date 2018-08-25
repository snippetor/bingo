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
)

type LocalCallFunc func(args interface{}, reply interface{}) error

type Server struct {
	name     string
	appName  string
	serv     *server.Server
	callFunc LocalCallFunc
}

func (s *Server) OnCall(ctx context.Context, args []byte, reply []byte) error {
	return s.callFunc(args, reply)
}

func (s *Server) Listen(name, appName string, port int, etcdAddrs []string, callFunc LocalCallFunc) {
	s.name = name
	s.appName = appName
	s.callFunc = callFunc
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
		s.serv.RegisterFunctionName(appName, "onCall", s.OnCall, "")
		if err := s.serv.Serve("kcp", ":"+strconv.Itoa(port)); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) Close() {
	s.serv.Close()
}

type Client struct {
	name     string
	appName  string
	addr     string
	client   client.XClient
	callFunc LocalCallFunc
}

func (c *Client) Connect(name, appName, serverAppName string, etcdAddrs []string) {
	c.name = name
	c.appName = appName
	d := client.NewEtcdDiscovery("bingo", serverAppName, etcdAddrs, nil)
	c.client = client.NewXClient(serverAppName, client.Failtry, client.RoundRobin, d, client.DefaultOption)
}

func (c *Client) Call(method string, args interface{}, reply interface{}) error {
	return c.client.Call(context.Background(), method, args, reply)
}

func (c *Client) CallNoReturn(method string, args interface{}) error {
	return c.client.Call(context.Background(), method, args, nil)
}

func (c *Client) CallNoReturn(method string, args interface{}) error {
	return c.client.Call(context.Background(), method, args, nil)
}

func (c *Client) Close() {
	c.client.Close()
}
