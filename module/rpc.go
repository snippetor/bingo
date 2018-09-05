package module

import (
	"github.com/snippetor/bingo/rpc"
)

// rpc
type RPCModule interface {
	Module
	Call(appName, method string, args interface{}, reply interface{}) error
}

type rpcModule struct {
	appName    string
	rpcClients []*rpc.Client
	rpcServer  *rpc.Server
}

func NewRPCModule(appName string, rpcClients []*rpc.Client, rpcServer *rpc.Server) RPCModule {
	return &rpcModule{appName, rpcClients, rpcServer}
}

func (m *rpcModule) Call(appName, method string, args interface{}, reply interface{}) error {
	if m.rpcClients != nil && len(m.rpcClients) > 0 {
		for _, c := range m.rpcClients {
			if c.ServerPkg == appName {
				return c.Call(method, args, reply)
			}
		}
	}
	return nil
}

func (m *rpcModule) CallNoReturn(appName, method string, args interface{}) error {
	return m.Call(appName, method, args, nil)
}

func (m *rpcModule) Close() {
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
