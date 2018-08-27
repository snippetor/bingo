package net

import (
	"net"
	"github.com/snippetor/bingo/log/fwlogger"
)

type tcpConn struct {
	baseConn
	conn *net.TCPConn
}

func (c *tcpConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(globalPacker.Pack(msgId, body))
		return true
	} else {
		fwlogger.W("-- send message failed!!! -- %#X", msgId)
		return false
	}
}

func (c *tcpConn) read(buf *[]byte) (int, error) {
	if c.conn != nil {
		return c.conn.Read(*buf)
	}
	return -1, nil
}

func (c *tcpConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *tcpConn) Address() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return "0:0:0:0"
}

func (c *tcpConn) NetProtocol() Protocol {
	return Tcp
}
