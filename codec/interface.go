package codec

import "github.com/snippetor/bingo/net"

// 数据传输协议，即包体格式
type ICodec interface {
	Marshal(interface{}) (net.MessageBody, error)
	Unmarshal(net.MessageBody, interface{}) error
}
