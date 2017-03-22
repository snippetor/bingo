package net

import (
	"github.com/snippetor/bingo"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"github.com/snippetor/bingo/comm"
	"sync"
)

type wsConn struct {
	conn *websocket.Conn
	absConn
}

func (c *wsConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.WriteMessage(websocket.BinaryMessage, GetDefaultMessagePacker().Pack(msgId, body))
		return true
	} else {
		bingo.W("-- send message failed!!! --")
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
	comm.Configable
	sync.RWMutex
	upgrader *websocket.Upgrader
	callback IMessageCallback
	clients  map[Identity]IConn
}

func (s *wsServer) wsHttpHandle(w http.ResponseWriter, r *http.Request) {
	if conn, err := s.upgrader.Upgrade(w, r, nil); err == nil {
		c := IConn(&wsConn{conn: conn})
		c.setState(STATE_CONNECTED)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		go s.handleConnection(c, s.callback)
	} else {
		bingo.E("-- ws build connection failed!!! --")
	}
}

func (s *wsServer) listen(port int, callback IMessageCallback) bool {
	if s.upgrader == nil {
		s.upgrader = &websocket.Upgrader{}
	}
	s.callback = callback
	s.clients = make(map[Identity]IConn, 0)
	http.HandleFunc("/", s.wsHttpHandle)
	if err := http.ListenAndServe("localhost:"+strconv.Itoa(port), nil); err != nil {
		bingo.E(err.Error())
		return false
	}
	return true
}

func (s *wsServer) Close() {
}

func (s *wsServer) handleConnection(conn IConn, callback IMessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			conn.setState(STATE_CLOSED)
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			s.Lock()
			delete(s.clients, conn.Identity())
			s.Unlock()
			break
		}
		packer := GetDefaultMessagePacker()
		id, body, _ := packer.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (s *wsServer) GetConnection(identity Identity) (IConn, bool) {
	s.RLock()
	defer s.RUnlock()
	if s.clients == nil {
		return nil, false
	} else {
		identity, ok := s.clients[identity]
		return identity, ok
	}
}

type wsClient struct {
	comm.Configable
	sync.Mutex
	conn IConn
}

func (c *wsClient) connect(serverAddr string, callback IMessageCallback) bool {
	c.conn.setState(STATE_CONNECTING)
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	bingo.I("Ws connect server ok :%s", serverAddr)
	if err != nil {
		bingo.E(err.Error())
		c.conn.setState(STATE_CLOSED)
		return false
	}
	c.conn = IConn(&wsConn{conn:conn})
	c.conn.setState(STATE_CONNECTED)
	c.handleConnection(c.conn, callback)
	return true
}

func (c *wsClient) handleConnection(conn IConn, callback IMessageCallback) {
	var buf []byte
	defer conn.Close()
	for {
		_, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			c.conn = nil
			break
		}
		packer := GetDefaultMessagePacker()
		id, body, _ := packer.Unpack(buf)
		if body != nil {
			callback(conn, id, body)
		}
	}
}

func (c *wsClient) Send(msgId MessageId, body MessageBody) bool {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.GetState() == STATE_CONNECTED {
		return c.conn.Send(msgId, body)
	} else {
		bingo.W("-- send tcp message failed!!! conn wrong state --")
	}
	return false
}

func (c *wsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
