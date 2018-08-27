package module

import "github.com/snippetor/bingo/net"

type Services map[string]net.Server

// service
type ServiceModule interface {
	Module
	GetService(name string) (net.Server, bool)
}

type serviceModule struct {
	servers Services
}

func NewServiceModule(servers Services) ServiceModule {
	return &serviceModule{servers}
}

func (m *serviceModule) GetService(name string) (net.Server, bool) {
	if m.servers != nil {
		if c, ok := m.servers[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

func (m *serviceModule) Close() {
	if m.servers != nil {
		for _, v := range m.servers {
			if v != nil {
				v.Close()
			}
		}
	}
}
