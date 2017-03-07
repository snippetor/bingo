package net

import (
	"net"
	"github.com/snippetor/bingo"
	"github.com/gorilla/websocket"
)

type wsServer struct {
	listener *net.TCPListener
}

func (s *wsServer) Listen(port int, callback iMessageCallback) bool {
	websocket.Upgrade()
}

func (s *wsServer) handleConnection(conn iConn, callback iMessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			bingo.E(err.Error())
			conn.Close()
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			break
		}
		byteBuffer = append(byteBuffer, buf[:l]...)
		packer := GetMessagePacker()
		for {
			id, body, remains := packer.unpack(byteBuffer)
			if body == nil || len(remains) == 0 {
				break
			} else {
				callback(conn, id, body)
			}
		}
	}
}

func (s *wsServer) Close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}
