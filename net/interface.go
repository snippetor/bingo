package net

type MessageId int32
type MessageBody []byte

const (
	MSGID_CONNECT_DISCONNECT = -1
)

type iMessagePacker interface {
	// 封包，传入消息ID和包体，返回字节集
	pack(MessageId, MessageBody) []byte
	// 解包，传入符合包结构的字节集，返回消息ID，包体，剩余内容
	unpack([]byte) (MessageId, MessageBody, []byte)
}

type IConn interface {
	Send(msgId MessageId, body MessageBody) bool
	read([]byte) (int, error)
	Close()
	Address() string
}

type iMessageCallback func(conn IConn, msgId MessageId, body MessageBody)

type iServer interface {
	listen(int, iMessageCallback)
	close()
}

type absServer struct {
	config map[string]string
}

func (s *absServer) setConfig(key, value string) {
	if s.config == nil {
		s.config = make(map[string]string)
	}
	s.config[key] = value
}

func (s *absServer) getConfig(key string) (string, bool) {
	v, ok := s.config[key]
	return v, ok
}
