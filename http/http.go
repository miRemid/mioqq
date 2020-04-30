package http

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/miRemid/mio"
	"github.com/miRemid/mioqq"
	"github.com/miRemid/mioqq/http/config"
)

const (
	// NoticeHandler = =
	NoticeHandler = iota
	// RequestHandler 请求处理函数
	RequestHandler
)

const (
	message = "message"
	notice  = "notice"
	request = "request"
)

const (
	// PerPrivate 私人消息
	PerPrivate = 1
	// PerGroup 群组消息
	PerGroup = 1 << 1
	// PerDiscuss 讨论组消息
	PerDiscuss = 1 << 2
)

// Server Http服务
type Server struct {
	engine *mio.Engine
	api    *mioqq.API
	secret string

	routers        map[string]PluginRouter
	noticeHandler  HandleFunc
	requestHandler HandleFunc

	logger *log.Logger
	Logger bool

	// Database
	db         *gorm.DB
	Driver     string
	ConnectURI string

	messageChan chan *CQContext
}

// New 新建HTTP
func New(api string) (*Server, error) {
	var client Server
	tmp, err := mioqq.New(config.TOKEN, api)
	if err != nil {
		return nil, err
	}
	client.api = tmp
	client.engine = mio.New()
	client.secret = config.SECRET
	client.routers = make(map[string]PluginRouter)
	client.noticeHandler = ignore
	client.requestHandler = ignore
	client.Logger = true
	client.logger = &log.Logger{}
	config.SetLogger(client.logger)

	client.messageChan = make(chan *CQContext, 0)
	return &client, nil
}

func (server *Server) receive(ctx *mio.Context) {
	// 接受信息
	// 1. 返回相应
	// 2. 转化为CQContext
	// 3. 加入任务队列
	var context CQContext
	context.Context = ctx
	context.API = server.api
	if err := ctx.ReadJSON(&context); err != nil {
		server.SendLog(Error, "解析数据失败: %v", err)
	} else {
		server.messageChan <- &context
		// go server.transport(&context)
	}
	ctx.JSON(204, nil)
}

// Engine 获取http服务端
func (server *Server) Engine() *mio.Engine {
	return server.engine
}

func (server *Server) noticeHTTP(ctx *mio.Context) {

}
func (server *Server) requestHTTP(ctx *mio.Context) {

}

// Server 开启服务
func (server *Server) Server(addr string) error {
	server.engine.Use(server.signature())
	server.engine.POST("/", server.receive)
	server.engine.GET("/", server.receive)

	api := server.engine.NewGroup("/api")
	{
		api.GET("notice", server.noticeHTTP)
		api.GET("request", server.requestHTTP)
	}
	go server.handleMessage()
	return server.engine.Server(addr)
}

func (server *Server) handleMessage() {
	for {
		var ctx = <-server.messageChan
		go server.transport(ctx)
	}
}

func (server *Server) transport(ctx *CQContext) {
	// 转发
	switch ctx.PostType {
	case message:
		server.m(ctx)
		break
	case notice:
		server.noticeHandler(ctx)
		break
	case request:
		server.requestHandler(ctx)
		break
	}
}

// Plugin 添加一个插件
func (server *Server) Plugin(cmd string, per int, handlers ...HandleFunc) {
	router := PluginRouter{}
	private, group, discuss := server.checkPermission(per)
	router.private = private
	router.group = group
	router.discuss = discuss
	router.handlers = append([]HandleFunc{}, handlers...)
	if _, ok := server.routers[cmd]; ok {
		server.SendLog(Warn, "%s 已存在，以覆盖处理\n", cmd)
	}
	server.routers[cmd] = router
	server.SendLog(Info, "cmd: %s --> %d handlers", cmd, len(handlers))
}

// On set the handler function
func (server *Server) On(handler HandleFunc, flag int) {
	switch flag {
	case NoticeHandler:
		server.noticeHandler = handler
		break
	case RequestHandler:
		server.requestHandler = handler
		break
	default:
		server.SendLog(Error, "参数错误")
		break
	}
}

func (server *Server) m(ctx *CQContext) {
	cmd, params := ctx.CmdParser(ctx.RawMessage, config.CMD...)
	if plugin, ok := server.routers[cmd]; ok {
		ctx.Params = params
		flag := false
		switch ctx.MessageType {
		case "private":
			if plugin.private {
				server.SendLog(Info, "私人消息，响应%s插件", cmd)
				flag = true
			}
			break
		case "group":
			if plugin.group {
				server.SendLog(Info, "群组消息，响应%s插件", cmd)
				flag = true
			}
			break
		case "discuss":
			if plugin.discuss {
				server.SendLog(Info, "讨论族消息，响应%s插件", cmd)
				flag = true
			}
			break
		}
		if flag {
			ctx.handlers = append([]HandleFunc{}, plugin.handlers...)
			ctx.handlers[0](ctx)
		}
	}
}
