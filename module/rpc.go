package module

import (
	"github.com/snippetor/bingo/rpc"
	"strings"
)

// rpc
type RPCModule interface {
	Module
	GetEndStub(nodeName string) (rpc.IEndStub, bool)
	GetEndStubsWithPrefix(nodeNamePrefix string) []rpc.IEndStub
	GetEndStubsWithModelName(nodeModelName string) []rpc.IEndStub
	GetEndStubWithPrefixAndBalancingSeed(nodeNamePrefix string, balancingSeed uint32) (rpc.IEndStub, bool)
	GetEndStubWithModelNameAndBalancingSeed(nodeModelName string, balancingSeed uint32) (rpc.IEndStub, bool)
}

type rpcModule struct {
	appName    string
	rpcClients []*rpc.Client
	rpcServer  *rpc.Server
}

func NewRPCModule(appName string, rpcClients []*rpc.Client, rpcServer *rpc.Server) RPCModule {
	return &rpcModule{appName, rpcClients, rpcServer}
}

func (m *rpcModule) GetEndStub(nodeName string) (rpc.IEndStub, bool) {
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

func (m *rpcModule) GetEndStubsWithPrefix(nodeNamePrefix string) []rpc.IEndStub {
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

func (m *rpcModule) GetEndStubsWithModelName(nodeModelName string) []rpc.IEndStub {
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
func (m *rpcModule) GetEndStubWithPrefixAndBalancingSeed(nodeNamePrefix string, balancingSeed uint32) (rpc.IEndStub, bool) {
	stubs := m.GetEndStubsWithPrefix(nodeNamePrefix)
	l := len(stubs)
	if l == 0 {
		return nil, false
	} else if l == 1 {
		return stubs[0], true
	} else {
		return stubs[balancingSeed%uint32(l)], true
	}
}

// 通过节点模型名称和平衡因子获取stub
// 如果有多个stub，则通过取模算法（balancingSeed % (stubs size)）来决定使用哪个
func (m *rpcModule) GetEndStubWithModelNameAndBalancingSeed(nodeModelName string, balancingSeed uint32) (rpc.IEndStub, bool) {
	stubs := m.GetEndStubsWithModelName(nodeModelName)
	l := len(stubs)
	if l == 0 {
		return nil, false
	} else if l == 1 {
		return stubs[0], true
	} else {
		return stubs[balancingSeed%uint32(l)], true
	}
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
