package net

var (
	globalPacker       iMessagePacker
	msgProtocol        iProtocol
	msgProtoCollection *protoCollection
)

func init() {
	globalPacker = iMessagePacker(&defaultMessagePacker{})
	msgProtoCollection = new(protoCollection)
}

func getMessagePacker() iMessagePacker {
	return globalPacker
}

type NetProtoType byte
type MessageProtoType byte

const (
	NetProtoType_TCP       NetProtoType = iota
	NetProtoType_WebSocket
)

const (
	MsgProto_Json     MessageProtoType = iota
	MsgProto_ProtoBuf
)

type ServerConfig struct {
	Net string
}

// 同步执行网络监听
// net: "tcp"/"ws"
func Listen(netProto NetProtoType, msgProto MessageProtoType, port int, callback IMessageCallback) bool {
	var server iServer
	switch netProto {
	case NetProtoType_TCP:
		server = iServer(&tcpServer{})
	case NetProtoType_WebSocket:
		server = iServer(&wsServer{})
	}
	switch msgProto {
	case MsgProto_Json:
		msgProtocol = iProtocol(&protocolJson{})
	case MsgProto_ProtoBuf:
		msgProtocol = iProtocol(&protocolProtoBuf{})
	}
	return server.listen(port, callback)
}

// 建立消息ID和消息结构的对应关系，用于数据解析
func PairMsgProto(msgId MessageId, msgStruct interface{}) {
	msgProtoCollection.put(msgId, msgStruct)
}

// 删除消息ID和消息结构的对应关系
func UnpairMsgProto(msgId MessageId) {
	msgProtoCollection.del(msgId)
}

// 清空消息ID和消息结构的对应关系
func ClearMsgProto(msgId MessageId) {
	msgProtoCollection.clear()
}
