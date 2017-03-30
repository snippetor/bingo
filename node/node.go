package node

import (
	"io/ioutil"
	"encoding/json"
	"github.com/snippetor/bingo"
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"strings"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/snippetor/bingo/codec"
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
	nodes  = make(map[string]*Model, 0)
)

func init() {
	config = &Config{}
}

func Parse(configPath string) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		bingo.E("-- parse config file failed! %s --", err)
		return
	}
	if err = json.Unmarshal(content, config); err != nil {
		bingo.E("-- parse config file failed! %s --", err)
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
		bingo.E("-- not found model by name %s --", n.ModelName)
		return
	}
	m.SetDefaultMessageProtocol(codec.Protobuf)
	m.SetClientProtoVersion("1.0")
	m.OnInit()
	// rpc
	if n.RpcPort > 0 {
		s := &rpc.Server{}
		s.Listen(n.RpcPort)
		m.rpcServer = s
	}
	for _, rpcServerNode := range n.RpcTo {
		c := &rpc.Client{}
		serverNode := findNode(rpcServerNode)
		c.Connect(config.Domains[serverNode.Domain])
		if m.rpcClients == nil {
			m.rpcClients = make(map[string]*rpc.Client)
		}
		m.rpcClients[n.Name] = c
	}
	// service
	for _, s := range n.Services {
		switch strings.ToLower(s.Type) {
		case "tcp":
			serv := net.GoListen(net.Tcp, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				if v, ok := m.proto.UnmarshalTo(msgId, body, m.clientProtoVer); ok {
					m.OnReceiveServiceMessage(conn, msgId, v)
				}
			})
			if m.servers == nil {
				m.servers = make(map[string]net.IServer)
			}
			m.servers[s.Name] = serv
		case "ws":
			serv := net.GoListen(net.WebSocket, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				if v, ok := m.proto.UnmarshalTo(msgId, body, m.clientProtoVer); ok {
					m.OnReceiveServiceMessage(conn, msgId, v)
				}
			})
			if m.servers == nil {
				m.servers = make(map[string]net.IServer)
			}
			m.servers[s.Name] = serv
		case "http":
			if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(ctx *fasthttp.RequestCtx) {
				m.OnReceiveHttpServiceRequest(ctx)
			}); err != nil {
				bingo.E("-- startup http service failed! %s --", err)
			}
		}
	}

	nodes[n.Name] = &m
}

func Run(nodeName string) {
	n := findNode(nodeName)
	if n == nil {
		bingo.E("-- run node failed! not found node by name %s --", nodeName)
		return
	}
	run_node(n)
}

func RunAll() {
	for _, n := range config.Nodes {
		run_node(n)
	}
}
