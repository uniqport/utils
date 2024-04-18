package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/uniqport/utils/pkg/client"
)

type uTcpClient struct {
	id        string
	tcpClient client.TcpClient
}

type TcpRoom struct {
	Server  *TcpServer
	roomId  string
	clients []uTcpClient
}

func NewTcpRoom(roomId string, address string) *TcpRoom {
	return &TcpRoom{
		roomId: roomId,
		Server: NewTcpServer(address),
	}
}

func (r *TcpRoom) Start() error {
	ln, err := net.Listen("tcp", r.Server.Address)
	if err != nil {
		return err
	}
	r.Server.listener = ln
	defer ln.Close()

	go r.acceptLoop()

	<-r.Server.quitch
	close(r.Server.msgch)

	return nil
}

func (r *TcpRoom) acceptLoop() {
	for {
		con, err := r.Server.listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Printf("New connection at %s \n", con.RemoteAddr().String())

		go r.handleConnection(con)

	}
}

func (r *TcpRoom) handleConnection(con net.Conn) {
	defer con.Close()
	var size int64
	binary.Read(con, binary.LittleEndian, &size)
	buf := new(bytes.Buffer)
	// TODO: add connection as a new room client
	for {
		n, err := io.CopyN(buf, con, size)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("recievd %d bytes over  \n", n)

		r.Server.msgch <- TcpMessage{
			payload: buf.Bytes(),
			from:    con.RemoteAddr().String(),
		}
	}
}

func (r *TcpRoom) HandleMessages() {
	for msg := range r.Server.msgch {
		fmt.Printf("[Room %s] | from: %s recieved: %s \n", r.roomId, msg.from, string(msg.payload))
	}
}
