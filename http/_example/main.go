package main

import (
	"github.com/miRemid/mioqq/http/plugins"

	"github.com/miRemid/mioqq/http"
)

func main() {
	server, err := http.New("http://127.0.0.1:5700")
	if err != nil {
		server.SendLog(http.Error, err.Error())
		return
	}
	server.Register(plugins.Roll{
		Cmd:  "roll",
		Area: 1000,
	})
	server.Server(":5678")
}
