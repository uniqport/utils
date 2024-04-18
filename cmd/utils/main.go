package main

import "github.com/uniqport/utils/pkg/server"

func main() {
	server := server.TcpServer(":3000")
	go server.HandleMessages()
	server.Start()
}
