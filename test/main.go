package main

import (
	"fmt"
	"os"
	"github.com/snippetor/bingo/net"
)

var ()

func main() {
	fmt.Println(os.Args)

	net.Listen(net.Kcp, 8080, nil)
}
