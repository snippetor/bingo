package net

import (
	"github.com/snippetor/bingo"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type wsConn struct {
	conn *websocket.Conn
}

func (c *wsConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.WriteMessage(websocket.BinaryMessage, GetDefaultMessagePacker().Pack(msgId, body))
		return true
	} else {
		bingo.E("-- send message failed!!! --")
		return false
	}
}

func (c *wsConn) read(buf *[]byte) (int, error) {
	if c.conn != nil {
		t, msg, err := c.conn.ReadMessage()
		if err == nil && t == websocket.BinaryMessage {
			*buf = msg
			return len(msg), nil
		}
	}
	return -1, nil
}

func (c *wsConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *wsConn) Address() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return "0:0:0:0"
}

func (c *wsConn) GetNetProtocol() NetProtocol {
	return WebSocket
}

type wsServer struct {
	absServer
	upgrader *websocket.Upgrader
	callback IMessageCallback
}

func (s *wsServer) wsHttpHandle(w http.ResponseWriter, r *http.Request) {
	if conn, err := s.upgrader.Upgrade(w, r, nil); err == nil {
		go s.handleConnection(IConn(&wsConn{conn: conn}), s.callback)
	} else {
		bingo.E("-- ws build connection failed!!! --")
	}
}

func (s *wsServer) listen(port int, callback IMessageCallback) bool {
	if s.upgrader == nil {
		s.upgrader = &websocket.Upgrader{}
	}
	s.callback = callback
	http.HandleFunc("/", s.wsHttpHandle)
	if err := http.ListenAndServe("localhost:"+strconv.Itoa(port), nil); err != nil {
		bingo.E(err.Error())
		return false
	}
	return true
}

func (s *wsServer) handleConnection(conn IConn, callback IMessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			break
		}
		packer := GetDefaultMessagePacker()
		id, body, _ := packer.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (s *wsServer) close() {
}
