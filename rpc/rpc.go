package rpc

import "github.com/snippetor/bingo/net"

type server struct {
	port int
}

func RunAsServer(port int) {
	go net.Listen(net.Tcp, port, handleMessage)
}

func RunAsClient(serverAddr string) {
	go net.Connect(net.Tcp, serverAddr, handleMessage)
}

func handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {

}
