package net

import (
	"net"
	"github.com/snippetor/bingo"
	"strconv"
	"github.com/snippetor/bingo/comm"
	"sync"
	"github.com/snippetor/bingo/utils"
)

type tcpConn struct {
	conn *net.TCPConn
	absConn
}

func (c *tcpConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(GetDefaultMessagePacker().Pack(msgId, body))
		return true
	} else {
		bingo.W("-- send message failed!!! --")
		return false
	}
}

func (c *tcpConn) read(buf *[]byte) (int, error) {
	if c.conn != nil {
		return c.conn.Read(*buf)
	}
	return -1, nil
}

func (c *tcpConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *tcpConn) Address() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	}
	return "0:0:0:0"
}

func (c *tcpConn) GetNetProtocol() NetProtocol {
	return Tcp
}

type tcpServer struct {
	comm.Configable
	sync.RWMutex
	listener *net.TCPListener
	clients  map[utils.Identity]IConn
}

func (s *tcpServer) listen(port int, callback IMessageCallback) bool {
	addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		bingo.E(err.Error())
		return false
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		bingo.E(err.Error())
		return false
	}
	defer listener.Close()
	s.listener = listener
	s.clients = make(map[utils.Identity]IConn, 0)
	bingo.I("Tcp server runnning on :%d", port)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		bingo.I(conn.RemoteAddr().String(), " tcp connect success")
		c := IConn(&tcpConn{conn: conn})
		c.setState(STATE_CONNECTED)
		s.Lock()
		s.clients[c.Identity()] = c
		s.Unlock()
		go s.handleConnection(c, callback)
	}
	return true
}

// 处理消息流
func (s *tcpServer) handleConnection(conn IConn, callback IMessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			conn.setState(STATE_CLOSED)
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			s.Lock()
			delete(s.clients, conn.Identity())
			s.Unlock()
			break
		}
		byteBuffer = append(byteBuffer, buf[:l]...)
		packer := GetDefaultMessagePacker()
		for {
			id, body, remains := packer.Unpack(byteBuffer)
			if body == nil || len(remains) == 0 {
				break
			} else {
				callback(conn, id, body)
			}
		}
	}
}

func (s *tcpServer) GetConnection(identity utils.Identity) (IConn, bool) {
	s.RLock()
	defer s.RUnlock()
	if s.clients == nil {
		return nil, false
	} else {
		identity, ok := s.clients[identity]
		return identity, ok
	}
}

func (s *tcpServer) Close() {
	s.Lock()
	defer s.Unlock()
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

type tcpClient struct {
	comm.Configable
	sync.Mutex
	conn IConn
}

func (c *tcpClient) connect(serverAddr string, callback IMessageCallback) bool {
	c.conn.setState(STATE_CONNECTING)
	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		bingo.E(err.Error())
		c.conn.setState(STATE_CLOSED)
		return false
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		bingo.E(err.Error())
		c.conn.setState(STATE_CLOSED)
		return false
	}
	defer conn.Close()
	c.conn = IConn(&tcpConn{conn: conn})
	c.conn.setState(STATE_CONNECTED)
	bingo.I("Tcp connect server ok :%s", serverAddr)
	c.handleConnection(c.conn, callback)
	return true
}

// 处理消息流
func (c *tcpClient) handleConnection(conn IConn, callback IMessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			c.conn.setState(STATE_CLOSED)
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
			c.conn = nil
			break
		}
		byteBuffer = append(byteBuffer, buf[:l]...)
		packer := GetDefaultMessagePacker()
		for {
			id, body, remains := packer.Unpack(byteBuffer)
			if body == nil || len(remains) == 0 {
				break
			} else {
				callback(conn, id, body)
			}
		}
	}
}

func (c *tcpClient) Send(msgId MessageId, body MessageBody) bool {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil && c.conn.GetState() == STATE_CONNECTED {
		return c.conn.Send(msgId, body)
	} else {
		bingo.W("-- send tcp message failed!!! conn wrong state --")
	}
	return false
}

func (c *tcpClient) Close() {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
