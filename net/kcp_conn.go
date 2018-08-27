package net

import (
	"github.com/snippetor/bingo/log/fwlogger"
	"net"
)

type kcpConn struct {
	baseConn
	conn net.Conn
}

func (c *kcpConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(globalPacker.Pack(msgId, body))
		return true
	} else {
		fwlogger.W("-- send message failed!!! -- %#X", msgId)
		return false
	}
}

func (c *kcpConn) read(buf *[]byte) (int, error) {
	if c.conn != nil {
		return c.conn.Read(*buf)
	}
	return -1, nil
}

func (c *kcpConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *kcpConn) Address() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return "0:0:0:0"
}

func (c *kcpConn) NetProtocol() Protocol {
	return Kcp
}
