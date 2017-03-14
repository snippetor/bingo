package net

import (
	"net"
	"github.com/snippetor/bingo"
	"strconv"
	"github.com/snippetor/bingo/comm"
)

type tcpConn struct {
	identity Identity
	conn     *net.TCPConn
}

func (c *tcpConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(GetDefaultMessagePacker().Pack(msgId, body))
		return true
	} else {
		bingo.E("-- send message failed!!! --")
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

func (c *tcpConn) Identity() Identity {
	if !isValidIdentity(c.identity) {
		c.identity = genIdentity()
	}
	return c.identity
}

type tcpServer struct {
	comm.Configable
	listener *net.TCPListener
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
	bingo.I("Tcp server runnning on :%d", port)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		bingo.I(conn.RemoteAddr().String(), " tcp connect success")
		go tcp_handleConnection(IConn(&tcpConn{conn}), callback)
	}
	return true
}

func (s*tcpServer) close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

type tcpClient struct {
	comm.Configable
	conn *net.TCPConn
}

func (c *tcpClient) connect(serverAddr string, callback IMessageCallback) bool {
	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		bingo.E(err.Error())
		return false
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		bingo.E(err.Error())
		return false
	}
	defer conn.Close()
	c.conn = conn
	bingo.I("Tcp connect server ok :%s", serverAddr)
	tcp_handleConnection(IConn(&tcpConn{conn}), callback)
	return true
}

func (c *tcpClient) close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// 处理消息流
func tcp_handleConnection(conn IConn, callback IMessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	defer conn.Close()
	for {
		l, err := conn.read(&buf)
		if err != nil {
			bingo.E(err.Error())
			callback(conn, MSGID_CONNECT_DISCONNECT, nil)
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
