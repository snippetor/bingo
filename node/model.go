package node

import (
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/proto"
	"github.com/snippetor/bingo/codec"
)

var (
	models = make(map[string]Model)
)

// rpc
// service
type Model struct {
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

func (m *Model) RegisterRPCMethod(methodName string, methd rpc.RPCMethod) {
	rpc.RegisterMethod(methodName, methd)
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

func BindNodeModel(modelName string, m Model) {
	models[modelName] = m
}

func getNodeModel(modelName string) (Model, bool) {
	m, ok := models[modelName]
	return m, ok
}
