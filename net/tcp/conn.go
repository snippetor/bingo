package tcp

import (
	"net"
)

type TCPConn struct {
	Address string
	Conn    *net.TCPConn
}

func (c *TCPConn) Forward(msg *ISCMessage) bool {
	return c.Send(msg.Id, msg.Body);
}

func (c *TCPConn) Send(msgId int32, msg []byte) bool {
	if c.Conn != nil {
		c.Conn.Write(net.Pack(msgId, msg))
		return true
	}
	return false
}

func (c *ISCConn) Read(b []byte) (int, error) {
	return this.Conn.Read(b)
}

func (this *ISCConn) Close() {
	if this.Conn != nil {
		this.Conn.Close()
		this.Conn = nil
	}
}