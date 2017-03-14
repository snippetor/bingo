package rpc

import "github.com/snippetor/bingo/net"

type Context struct {
	conn    net.IConn
	method  string
	args    map[string]string
	version string
}
