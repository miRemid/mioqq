package http

import (
	"log"

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

// Server Http服务
type Server struct {
	engine *mio.Engine
	api    *mioqq.API
	secret string

	plugins        map[string]p
	noticeHandler  HandleFunc
	requestHandler HandleFunc

	logger *log.Logger
	Logger bool
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

	client.plugins = make(map[string]p)
	client.noticeHandler = ignore
	client.requestHandler = ignore

	client.Logger = true
	client.logger = &log.Logger{}
	config.SetLogger(client.logger)

	return &client, nil
}

func (server *Server) receive(ctx *mio.Context) {
	var context CQContext
	context.Context = ctx
	context.API = server.api
	if err := ctx.ReadJSON(&context); err != nil {
		server.SendLog(Error, "解析数据失败: %v", err)
		ctx.JSON(204, nil)
		return
	}
	server.transport(&context)
}

// Engine 获取http服务端
func (server *Server) Engine() *mio.Engine {
	return server.engine
}

// Server 开启服务
func (server *Server) Server(addr string) error {
	server.engine.Use(server.signature())
	server.engine.POST("/", server.receive)
	server.engine.GET("/", server.receive)
	return server.engine.Server(addr)
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
	ctx.JSON(204, nil)
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
	cmd, _ := ctx.CmdParser(ctx.RawMessage, config.CMD...)
	if plugin, ok := server.plugins[cmd]; ok {
		flag := false
		switch ctx.MessageType {
		case "private":
			if plugin.private {
				server.SendLog(Info, "私人消息，响应%s插件", plugin.name)
				flag = true
			}
			break
		case "group":
			if plugin.group {
				server.SendLog(Info, "群组消息，响应%s插件", plugin.name)
				flag = true
			}
			break
		case "discuss":
			if plugin.discuss {
				server.SendLog(Info, "讨论族消息，响应%s插件", plugin.name)
				flag = true
			}
			break
		}
		if flag {
			plugin.plugin.Parse(ctx)
		}
	}
}
