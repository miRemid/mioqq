# MioQQ
mioqq是一个golang编写的CQHTTP SDK&http服务

# 快速使用

```golang

package main

import (
    "log"
    "github.com/miRemid/mioqq"
)

func main() {
    // 填写token和CQHTTP API
    m, _ := mioqq.New("", "http://127.0.0.1:5700")

    // qq号，消息类型，消息格式
    message := m.NewMessage(
        123456, 
        mioqq.PrivateMessage, 
        mioqq.StringContent
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
-   mioqq.PrivateMessage
-   mioqq.GroupMessage
-   mioqq.DiscussMessage  

format指明发送消息的消息格式，可选参数为
-   mioqq.StringContent(字符串格式)
-   mioqq.ArrayContent(消息段格式)

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

# HTTP服务
`mioqq`已经提供了一个http服务位于`mioqq/http`中

## 快速使用
```go
package main

import (
	"github.com/miRemid/mioqq/http/plugins"

	"github.com/miRemid/mioqq/http"
)

func main() {
	server, err := http.New("CQHTTP_API")
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
```
`mioqq/http`需要一个配置文件，其字段如下
```json
{   
    // 应用名称
    "name": "mio",
    // 日志目录，默认当前目录下的log
    "log_path": "./log",
    // 签名密钥，默认为空
    "screct": "mio",
    // access_token，api请求token
    "access_token": "asdf",
    // cmd 命令标识符号
    "cmd": ["!", "！", "###"]
}
```
运行之后，对目标QQ输入`!roll`即可触发插件并返回一个随机数
> 服务器默认对好友请求或其他请求消息、Notice消息不以处理

## Plugin插件接口
`mioqq/http`的插件需要满足`mioqq/http/Plugin`接口，其结构如下
```golang
type Plugin interface{
    Parse(ctx *http.CQContext)
    Help() string
}
```
## Plugin标准结构示范
一个标准的`mioqq/http`插件必须满足以下条件:
1. 满足`Plugin`接口
2. 结构体必须包含`Cmd`字段，并且该字段的`tag`中可以需要包含`mio`和`role`字段
>Cmd的tag字段中，mio的作用是指明命令处罚条件如：roll、hello等，如果添加了该字段，则在Register的时候无需修改Cmd的值
例如以下就是一个标准的`Plugin`
```golang
type Example struct {
    Cmd string `mio:"haha" role:"7"`
}

func (e Example) Parse(ctx *http.CQContext) {
    if res, err := evt.Send("hahaha", true, false); err != nil {
        ...
    }else {
        ...
    }
}

func (e Example) Help() string {
    return "return hahaha"
}
```
## Plugin 响应消息权限
默认所有命令响应群组、私人、讨论组消息。如需更改权限，在插件的Cmd字段中添加tag:`role`.
```golang
type Example struct {
    Cmd string `mio:"example" role:"7"`
}
```
hanabi将会提取role字段转为int类型取低三位数字，根据其二进制判断消息权限

    000 不响应任何消息
    001 私人消息，1
    010 群组消息，2
    100 讨论祖消息，4
    111 所有消息，7