package client_test

import (
	"testing"

	"github.com/uniqport/utils/pkg/client"
	"github.com/uniqport/utils/pkg/server"
)

const (
	port     string = ":3000"
	username string = "testClient"
)

func TestTcpClient(t *testing.T) {
	testClient := client.TcpClient(port, username)

	if testClient.Address != port {
		t.Fatalf("expected port: %s, got: %s", port, testClient.Address)
	}

	if testClient.Username != username {
		t.Fatalf("expected port: %s, got: %s", username, testClient.Username)
	}
}

func TestDialAndSend(t *testing.T) {
	testServer := server.TcpServer(port)
	go testServer.Start()
	testClient := client.TcpClient(port, username)
	if err := testClient.Connect(); err != nil {
		t.Fatalf("error connecting to the server %v", err)
	}

	if err := testClient.Send("test"); err != nil {
		t.Fatalf("error sending messages: %v", err)
	}
}
