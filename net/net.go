package net

var (
	globalPacker iMessagePacker
)

func init() {
	globalPacker = iMessagePacker(&defaultMessagePacker{})
}

func getMessagePacker() iMessagePacker {
	return globalPacker
}

type NetProtocol byte
type MessageProtocol byte

const (

)

type ServerConfig struct {
	Net string
}

// 同步执行网络监听
// net: "tcp"/"ws"
func Listen(net string, port int, callback IMessageCallback) bool {
	var server iServer
	switch net {
	case "tcp":
		server = iServer(&tcpServer{})
	case "ws":
		server = iServer(&wsServer{})
	}
	return server.listen(port, callback)
}
