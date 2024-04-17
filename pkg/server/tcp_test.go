package server_test

import (
	"testing"

	"github.com/uniqport/utils/pkg/server"
)

const port string = ":3000"

func TestTcpServer(t *testing.T) {
	testServer := server.TcpServer(port)

	if testServer.Address != port {
		t.Fatalf("expected port: %s, got: %s", port, testServer.Address)
	}
}
