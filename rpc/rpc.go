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
	"github.com/smallnest/rpcx/server"
	"strconv"
	"github.com/smallnest/rpcx/client"
	"context"
	"github.com/smallnest/rpcx/protocol"
	"github.com/gogo/protobuf/proto"
	"errors"
	"github.com/snippetor/bingo/codec"
	"github.com/snippetor/bingo/net"
)

type Server struct {
	appName string
	pkg     string
	serv    *server.Server
}

func (s *Server) Listen(appName, pkg string, port int, onInit func(server *Server)) {
	s.appName = appName
	s.pkg = pkg
	s.serv = server.NewServer()
	go func() {
		if onInit != nil {
			onInit(s)
		}
		if err := s.serv.Serve("kcp", ":"+strconv.Itoa(port)); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) RegisterFunction(name string, fn interface{}) {
	s.serv.RegisterFunctionName(s.pkg, name, fn, "")
}

func (s *Server) Close() {
	if s.serv != nil {
		s.serv.Close()
	}
}

type Client struct {
	appName   string
	pkg       string
	addr      string
	ServerPkg string
	client    client.XClient
}

func (c *Client) Connect(appName, pkg, serverPkg string, address []string) {
	c.appName = appName
	c.pkg = pkg
	c.ServerPkg = serverPkg
	var pairs []*client.KVPair
	for i := range address {
		pairs = append(pairs, &client.KVPair{Key: address[i]})
	}
	d := client.NewMultipleServersDiscovery(pairs)
	var option = client.DefaultOption
	option.SerializeType = protocol.SerializeNone
	c.client = client.NewXClient(serverPkg, client.Failtry, client.RoundRobin, d, option)
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
	if c.client != nil {
		c.client.Close()
	}
}
