package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/miRemid/mio"
)

func ignore(ctx *CQContext) {
	ctx.JSON(204, nil)
}

// Signature 消息验证中间件
func (server *Server) signature() mio.HandlerFunc {
	server.SendLog(Info, "以开启Signature验证, key=%v\n", server.secret)
	return func(ctx *mio.Context) {
		if server.secret != "" {
			sig := ctx.Req.Header.Get("X-Signature")
			if sig == "" {
				server.SendLog(Warn, "未找到X-Signature头部信息，请检查CQHTTP配置")
				ctx.JSON(204, nil)
				return
			}
			sig = sig[len("sha1="):]
			mac := hmac.New(sha1.New, []byte(server.secret))
			byteData, _ := ioutil.ReadAll(ctx.Req.Body)
			io.WriteString(mac, string(byteData))
			res := fmt.Sprintf("%x", mac.Sum(nil))
			if res != sig {
				server.SendLog(Info, "消息不来自CQHTTP，以屏蔽处理")
				ctx.JSON(204, nil)
				return
			}
			// 重写数据
			ctx.Req.Body = ioutil.NopCloser(bytes.NewReader(byteData))
			server.SendLog(Info, "接受到CQHTTP消息，开始解析处理")
		}
		ctx.Next()
	}
}
