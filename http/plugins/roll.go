package plugins

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/miRemid/mioqq"

	"github.com/miRemid/mioqq/http"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Roll 随机数插件
type Roll struct {
	Cmd  string `mio:"cmd" role:"7"`
	Area int
}

// Parse 解析函数
func (r Roll) Parse(ctx *http.CQContext) {
	if r.Area == 0 {
		r.Area = 100
	}
	text := fmt.Sprintf("%d", rand.Intn(r.Area))
	message := ctx.API.NewMessage(ctx.UserID, mioqq.PrivateMessage, mioqq.StringContent)
	message.Text(text)
	_, err := ctx.API.Send(message)
	if err != nil {
		fmt.Println(err)
		return
	}
}
