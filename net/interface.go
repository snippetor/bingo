package net

type IMessagePacker interface {
	// 封包，传入消息ID和包体，返回字节集
	Pack(int32, []byte) []byte
	// 解包，传入符合包结构的字节集，返回消息ID，包体，剩余内容
	Unpack([]byte) (int32, []byte, []byte)
}

type IClient interface {
	Connect()
	Close()
}

type IServer interface {
}

type IConn interface {
}
