package net

import (
	"net"
	"github.com/snippetor/bingo"
	"strconv"
)

type tcpConn struct {
	conn *net.TCPConn
}

func (c *tcpConn) Send(msg []byte) bool {
	if c.conn != nil && msg != nil && len(msg) > 0 {
		c.conn.Write(msg)
		return true
	} else {
		bingo.E("-- send message failed!!! --")
		return false
	}
}

func (c *tcpConn) read(buf []byte) (int, error) {
	if c.conn != nil {
		return c.conn.Read(buf)
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
	listener *net.TCPListener
}

func (s *tcpServer) listen(port int, callback iMessageCallback) bool {
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
	bingo.I("Tcp server runnning on :%d", strconv.Itoa(port))
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

func (s *tcpServer) handleConnection(conn IConn, callback iMessageCallback) {
	buf := make([]byte, 4096) // 4KB
	byteBuffer := make([]byte, 0)
	for {
		l, err := conn.read(buf)
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

func (s*tcpServer) close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}
