package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/buger/jsonparser"
)

var (
	file []byte
	// LogPath 日志文件目录
	LogPath string
	// TOKEN CQHTTP TOKEN
	TOKEN string
	// SCRECT CQHTTP SCRECT
	SECRET string
	// CMD 命令解析
	CMD []string
	// NAME 服务名称
	NAME string
	// AllowAllRequest 允许所有加好友请求
	AllowAllRequest bool
)

func init() {
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	file = f
	LoadString(&LogPath, "./log", "log_path")
	LoadString(&SECRET, "", "screct")
	LoadString(&TOKEN, "", "access_token")
	LoadString(&NAME, "mio", "name")
	LoadStringArray(&CMD, []string{"!"}, "cmd")
}

// LoadString 读入string配置
func LoadString(value *string, def string, keys ...string) (string, error) {
	res, err := jsonparser.GetString(file, keys...)
	if value == nil {
		return res, err
	}
	if err != nil {
		*value = def
	} else {
		*value = res
	}
	return res, err
}

// LoadStringArray 读入string数组配置
func LoadStringArray(value *[]string, def []string, keys ...string) (res []string, err error) {
	data, _, _, e := jsonparser.Get(file, keys...)
	if e != nil {
		return res, e
	}
	err = json.Unmarshal(data, &res)
	if value == nil {
		return res, err
	}
	if err != nil {
		*value = def
	} else {
		*value = res
	}
	return
}

// LoadInt 读入int配置
func LoadInt(value *int64, def int64, keys ...string) (int64, error) {
	res, err := jsonparser.GetInt(file, keys...)
	if value == nil {
		return res, err
	}
	if err != nil {
		*value = def
	} else {
		*value = res
	}
	return res, err
}

// LoadBoolean 读入bool配置
func LoadBoolean(value *bool, def bool, keys ...string) (bool, error) {
	res, err := jsonparser.GetBoolean(file, keys...)
	if value == nil {
		return res, err
	}
	if err != nil {
		*value = def
	} else {
		*value = res
	}
	return res, err
}

// LoadFloat 读入float配置
func LoadFloat(value *float64, def float64, keys ...string) (float64, error) {
	res, err := jsonparser.GetFloat(file, keys...)
	if value == nil {
		return res, err
	}
	if err != nil {
		*value = def
	} else {
		*value = res
	}
	return res, err
}
