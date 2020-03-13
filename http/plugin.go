package http

import (
	"fmt"
	"reflect"
	"strconv"
)

// Plugin is the plugin of hanabi
type Plugin interface {
	Parse(ctx *CQContext)
	Help() string
}

// HandleFunc 处理函数
type HandleFunc func(*CQContext)

type p struct {
	plugin  Plugin
	name    string
	private bool
	group   bool
	discuss bool
}

func (server *Server) permision(v reflect.Value, f reflect.StructField) (per int) {
	role := f.Tag.Get("role")
	if role == "" {
		server.SendLog(Warn, "[WA] %s插件权限读取失败，以设置默认权限为7", v.Type())
		return 7
	}
	tmp, err := strconv.Atoi(role)
	if err != nil {
		server.SendLog(Warn, "[WA] %s插件权限读取失败，以设置默认权限为7", v.Type())
		per = 7
	} else {
		per = tmp
	}
	return per
}

func (server *Server) checkPermission(role int) (bool, bool, bool) {
	private := role & 1
	group := role >> 1 & 1
	discuss := role >> 2 & 1
	return private == 1, group == 1, discuss == 1
}

// Register a plugin
func (server *Server) Register(pluginss ...Plugin) {
	for _, plugin := range pluginss {
		var cmd string
		var role int
		v := reflect.ValueOf(plugin)
		t := reflect.TypeOf(plugin)
		if f, ok := t.FieldByName("Cmd"); !ok {
			server.SendLog(Error, "%s插件读取失败，检查是否包含Cmd字段", v.Type())
			continue
		} else {
			cmd = fmt.Sprintf("%s", v.FieldByName("Cmd"))
			if cmd == "" {
				cmd = f.Tag.Get("hana")
			}
			role = server.permision(v, f)
		}
		if cmd == "" {
			server.SendLog(Error, "%s插件读取失败，检查初始化是否正确或tag是否包含hana字段", v.Type())
			continue
		}
		if _, ok := server.plugins[cmd]; ok {
			server.SendLog(Info, "%s插件覆盖成功", v.Type())
		} else {
			server.SendLog(Info, "%s插件读取成功", v.Type())
		}
		private, group, discuss := server.checkPermission(role)
		server.plugins[cmd] = p{
			plugin:  plugin,
			name:    fmt.Sprintf("%s", v.Type()),
			private: private,
			group:   group,
			discuss: discuss,
		}
	}
}
