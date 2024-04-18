package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type TcpClient struct {
	con      net.Conn
	Address  string
	Username string
}

func NewTcpClient(address string, username string) *TcpClient {
	return &TcpClient{
		Address:  address,
		Username: username,
	}
}

func (c *TcpClient) Connect() error {
	con, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	fmt.Printf("Connection established on: %s! \n", c.Address)

	c.con = con
	return nil
}

func (c *TcpClient) Send(msg string) error {
	if c.con == nil {
		return errors.New("there is no connection")
	}

	msgSize := len(msg)
	data := make([]byte, msgSize)
	binary.Write(c.con, binary.LittleEndian, int64(msgSize))

	n, err := io.CopyN(c.con, bytes.NewReader(data), int64(msgSize))
	if err != nil {
		return err
	}

	fmt.Printf("sent %d bytes over %s \n", n, c.Address)
	return nil
}
