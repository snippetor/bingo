package net

import (
	"net"
	"github.com/snippetor/bingo/errors"
)

type kcpConn struct {
	baseConn
	conn net.Conn
}

func (c *kcpConn) Send(msgId MessageId, body MessageBody) error {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(globalPacker.Pack(msgId, body))
		return nil
	} else {
		return errors.ConnectionError(errors.ErrCodeConnect)
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
