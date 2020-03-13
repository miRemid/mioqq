# Mio

# Install

    go get github.com/mio-qq/mio

# Quick Start

```golang

package main

import (
    "log"
    "github.com/mio-qq/mio"
)

func main() {
    // Token, CQHTTP接口, Screct
    // Token和Screct没有则填""
    m, _ := mio.NewMio("", "http://127.0.0.1:5700", "mio")

    // qq号，消息类型，消息格式
    message := m.NewMessage(
        123456, 
        mio.PrivateMessage, 
        mio.StringContent
    ).Text("Hello").NewLine().Face(1)

    data, _ := message.Parse()
    log.Println(string(data))
    res, err := m.Send(message)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(string(res.Data))
}
```

# Message
`Mio`发送的每一条消息都是`Mio.Message`结构，通过`mio.NewMessage`可以生成一条消息
```go
msg := mio.NewMessage(id, type, format)
```
id指明发送方的id号，如QQ号，群号，讨论组编号(COOLQ提供)  
type指发送消息类型，可选参数为:
-   mio.PrivateMessage
-   mio.GroupMessage
-   mio.DiscussMessage  

format指明发送消息的消息格式，可选参数为
-   mio.StringContent(字符串格式)
-   mio.ArrayContent(消息段格式)

将消息封装成`Message`的好处是可以方便生成消息，其添加内容的格式为`msg.FuncName(params)`，每个函数都将返回一个新的`Message`因此你可以链式生成一条消息，如:
```golang
msg.Text("hello").NewLine().Face(1)
```
以下是目前`Message`支持的常用消息API:
    
    NewLine(还行)
    Text(文本:string)
    Face(表情id:int)
    Emoji(表情id:int)
    At(用户id:int)
    AtAll(群组或讨论组可用)
    ...

# Server
`mio`自带了服务端服务`mio.Server`，轻量级处理了CQHTTP上报的消息，方便业务代码开发  
目前`Server`支持的服务端类型有`HTTP`和`Websocket`，`Server`将会自动根据用户提供的ip地址建立合适的服务

## Quick Start
```go
package main

import (
    "github.com/mio-qq/mio"
)

func messageHandler(session CQSession) {
	log.Println(session.Message)
	msg := session.NewMessage(351968703, PrivateMessage, ArrayContent)
	msg.Text("Websocket Teest")
	_, err := session.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func main() {
    // 创建服务
    server := NewServer()
    // 开启Debug模式，默认关闭
    server.DebugMode = true
    // 配置CQHTTP API信息
    // Websocket API
    // server.Config("", "ws://127.0.0.1:6700")
    // HTTP API
    server.Config("", "http://127.0.0.1:5700")
    // 设置消息处理函数
    server.MessageHandler = messageHandler
    // 如果服务端采用HTTP模式且CQHTTP配置了Screct
    server.Screct = "Your CQHTTP Config"

    // 开启服务
    // ws: ws://uri:port
    // http: http://uri:port or :port
	if err := server.Run("http://127.0.0.1:6700"); err != nil {
		log.Fatal(err)
	}
}
```