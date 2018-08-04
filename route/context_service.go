package route

import (
	"github.com/snippetor/bingo/net"
	"fmt"
)

type ServiceContext struct {
	Context
	Conn         net.IConn
	MessageType  int32
	MessageGroup int32
	MessageExtra int32
	MessageId    int32 // unpacked id
	MessageBody  net.MessageBody
}

// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*MyContext".
func (c *ServiceContext) Do(handlers Handlers) {
	Do(c, handlers)
}

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*MyContext".
func (c *ServiceContext) Next() {
	Next(c)
}

func (c *ServiceContext) Ack(body net.MessageBody) {
	if c.Conn == nil {
		panic(fmt.Sprintf("[ack] lost connection, app=%s ctx=%v", c.App().Name(), c.Id()))
	}
	c.Conn.Send(net.PackId(net.MsgTypeAck, c.MessageGroup, c.MessageExtra, c.MessageId), body)
}

//Note that: just send to the requester
func (c *ServiceContext) Ntf(body net.MessageBody) {
	if c.Conn == nil {
		panic(fmt.Sprintf("[ntf] lost connection, app=%s ctx=%v", c.App().Name(), c.Id()))
	}
	c.Conn.Send(net.PackId(net.MsgTypeNtf, c.MessageGroup, c.MessageExtra, c.MessageId), body)
}

//Note that: just send to the requester
func (c *ServiceContext) Cmd(body net.MessageBody) {
	if c.Conn == nil {
		panic(fmt.Sprintf("[cmd] lost connection, app=%s ctx=%v", c.App().Name(), c.Id()))
	}
	c.Conn.Send(net.PackId(net.MsgTypeCmd, c.MessageGroup, c.MessageExtra, c.MessageId), body)
}
