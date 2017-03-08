package main

import (
	"github.com/snippetor/bingo/net"
	"fmt"
)

func main() {
	net.Listen("tcp", 9090, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
		fmt.Println(msgId)
	})
}
