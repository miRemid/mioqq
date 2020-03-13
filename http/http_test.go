package http

import (
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	server, err := New("http://127.0.0.1:5700")
	if err != nil {
		log.Fatal(err)
	}
	server.Server(":5678")
}
