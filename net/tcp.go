package net

import (
	"net"
	"github.com/snippetor/bingo"
	"strconv"
)

type tcpConn struct {
	conn *net.TCPConn
}

func (c *tcpConn) Send(msgId MessageId, body MessageBody) bool {
	if c.conn != nil && body != nil && len(body) > 0 {
		c.conn.Write(getMessagePacker().pack(msgId, body))
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

type tcpServer struct {
	absServer
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
		go s.handleConnection(IConn(&tcpConn{conn}), callback)
	}
	return true
}

func (s *tcpServer) handleConnection(conn IConn, callback IMessageCallback) {
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
		packer := getMessagePacker()
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

func (s*tcpServer) close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}
