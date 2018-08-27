package app

import (
	"github.com/snippetor/bingo/codec"
	"github.com/snippetor/bingo/net"
)

type ServiceContext struct {
	Context
	Conn         net.Conn
	MessageType  int32
	MessageGroup int32
	MessageExtra int32
	MessageId    int32 // unpacked id
	MessageBody  *MessageBodyWrapper
	Codec        codec.ICodec
}

type MessageBodyWrapper struct {
	RawContent net.MessageBody
	Codec      codec.ICodec
}

func (c *MessageBodyWrapper) Unmarshal(v interface{}) {
	c.Codec.Unmarshal(c.RawContent, v)
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

func (c *ServiceContext) Ack(body interface{}) {
	if c.Conn == nil {
		c.LogE("[ack] lost connection, app=%s ctx=%v", c.App().Name(), c.Id())
		return
	}
	switch body.(type) {
	case []byte, net.MessageBody:
		c.Conn.Send(net.PackId(net.MsgTypeAck, c.MessageGroup, c.MessageExtra, c.MessageId), body.([]byte))
	default:
		c.Conn.Send(net.PackId(net.MsgTypeAck, c.MessageGroup, c.MessageExtra, c.MessageId), c.Codec.Marshal(body))
	}
	// ACK will stop the execution queue.
	c.StopExecution()
}

//Note that: just send to the requester
func (c *ServiceContext) Ntf(body interface{}) {
	if c.Conn == nil {
		c.LogE("[ntf] lost connection, app=%s ctx=%v", c.App().Name(), c.Id())
		return
	}
	switch body.(type) {
	case []byte, net.MessageBody:
		c.Conn.Send(net.PackId(net.MsgTypeNtf, c.MessageGroup, c.MessageExtra, c.MessageId), body.([]byte))
	default:
		c.Conn.Send(net.PackId(net.MsgTypeNtf, c.MessageGroup, c.MessageExtra, c.MessageId), c.Codec.Marshal(body))
	}
}

//Note that: just send to the requester
func (c *ServiceContext) Cmd(body interface{}) {
	if c.Conn == nil {
		c.LogE("[cmd] lost connection, app=%s ctx=%v", c.App().Name(), c.Id())
		return
	}
	switch body.(type) {
	case []byte, net.MessageBody:
		c.Conn.Send(net.PackId(net.MsgTypeCmd, c.MessageGroup, c.MessageExtra, c.MessageId), body.([]byte))
	default:
		c.Conn.Send(net.PackId(net.MsgTypeCmd, c.MessageGroup, c.MessageExtra, c.MessageId), c.Codec.Marshal(body))
	}
}
