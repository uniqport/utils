package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type Server interface {
	Start() error
	acceptLoop()
	handleConnection(net.Conn)
}

type TcpMessage struct {
	// TODO: get from info from user type
	from    string
	payload []byte
}

type TcpServer struct {
	quitch   chan int
	msgch    chan TcpMessage
	listener net.Listener
	Address  string
}

func NewTcpServer(address string) *TcpServer {
	return &TcpServer{
		Address: address,
		quitch:  make(chan int),
		msgch:   make(chan TcpMessage),
	}
}

func (s *TcpServer) Start() error {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	s.listener = ln
	defer ln.Close()

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *TcpServer) acceptLoop() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Printf("New connection at %s \n", con.RemoteAddr().String())

		go s.handleConnection(con)
	}
}

func (s *TcpServer) handleConnection(con net.Conn) {
	defer con.Close()
	var size int64
	binary.Read(con, binary.LittleEndian, &size)
	buf := new(bytes.Buffer)
	for {
		n, err := io.CopyN(buf, con, size)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("recievd %d bytes over  \n", n)

		s.msgch <- TcpMessage{
			payload: buf.Bytes(),
			from:    con.RemoteAddr().String(),
		}
	}
}

func (s *TcpServer) HandleMessages() {
	for msg := range s.msgch {
		fmt.Printf("from: %s recieved: %s \n", msg.from, string(msg.payload))
	}
}

// WARN: not sure if thats the way
func (s *TcpServer) Close() {
	s.quitch <- 0
}
