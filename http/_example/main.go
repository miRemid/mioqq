package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/miRemid/mioqq"

	"github.com/miRemid/mioqq/http"
)

func roll(ctx *http.CQContext) {
	text := fmt.Sprintf("%d", rand.Intn(1000))
	message := ctx.API.NewMessage(ctx.UserID, mioqq.PrivateMessage, mioqq.StringContent)
	message.Text(text)
	_, err := ctx.API.Send(message)
	if err != nil {
		fmt.Println(err)
		return
	}

}

type Help struct {
	Cmd string
}

func (help Help) Parse(ctx *http.CQContext) {
	_, err := ctx.Send("roll, 随机一个数字", false)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	server, err := http.New("http://127.0.0.1:5700")
	if err != nil {
		server.SendLog(http.Error, err.Error())
		return
	}
	server.Register(Help{
		Cmd: "help",
	})
	server.Plugin("roll", http.PerPrivate, roll)
	server.Server(":5678")
}
