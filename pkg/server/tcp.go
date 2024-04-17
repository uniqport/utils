package server

import (
	"fmt"
	"log"
	"net"
)

type tcpMessage struct {
	// TODO: get from info from user type
	from    string
	payload []byte
}

type tcpServer struct {
	quitch   chan struct{}
	msgch    chan tcpMessage
	listener net.Listener
	Address  string
}

func TcpServer(address string) *tcpServer {
	return &tcpServer{
		Address: address,
		quitch:  make(chan struct{}),
		msgch:   make(chan tcpMessage, 10),
	}
}

func (s *tcpServer) Start() error {
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

func (s *tcpServer) acceptLoop() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		s.handleConnection(con)
	}
}

func (s *tcpServer) handleConnection(con net.Conn) {
	// TODO: create the size of the value that will be transmitted
	buf := make([]byte, 2048)
	defer con.Close()
	for {
		n, err := con.Read(buf)
		if err != nil {
			log.Fatal(err)
			continue
		}

		s.msgch <- tcpMessage{
			payload: buf[:n],
			from:    con.RemoteAddr().String(),
		}
	}
}

func (s *tcpServer) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("from: %s recieved: %s \n", msg.from, string(msg.payload))
	}
}
