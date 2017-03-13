package rpc

import "github.com/snippetor/bingo/net"

type server struct {
	port int
}

func NewServer(port int) {
	go net.Listen(net.Tcp, port, handleMessage)
}

func handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {

}
