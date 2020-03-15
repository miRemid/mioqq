package http

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/miRemid/mio"
	"github.com/miRemid/mioqq"
)

// CQParams 参数
type CQParams map[string]interface{}

// CQContext 用户对话
type CQContext struct {
	Context  *mio.Context
	API      *mioqq.API
	handlers []HandleFunc
	index    int
	Params   []string

	quick bool
	ws    bool

	Time   int64 `json:"time"`
	SelfID int64 `json:"self_id"`

	PostType      string      `json:"post_type"`
	MessageType   string      `json:"message_type,omitempty"`
	SubType       string      `json:"sub_type,omitempty"`
	MessageID     int64       `json:"message_id,omitempty"`
	GroupID       int64       `json:"group_id,omitempty"`
	DiscussID     int64       `json:"discuss_id,omitempty"`
	UserID        int64       `json:"user_id"`
	Font          int32       `json:"font,omitempty"`
	Message       interface{} `json:"message,omitempty"`
	RawMessage    string      `json:"raw_message,omitempty"`
	Anonymous     interface{} `json:"anonymous,omitempty"`
	AnonymousFlag string      `json:"anonymous_flag,omitempty"`
	Sender        *User       `json:"sender"`

	NoticeType  string `json:"notice_type,omitempty"`
	OperatorID  int64  `json:"operator_id,omitempty"`
	File        *File  `json:"file,omitempty"`
	RequestType string `json:"request_type"`
	Flag        string `json:"flag,omitempty"`
	Comment     string `json:"comment"`
	Duration    int64  `json:"duration"`
}

// Next 进行下一个中间件
func (context *CQContext) Next() {
	context.index++
	s := len(context.handlers)
	if context.index != s {
		context.handlers[context.index](context)
	}
}

// CmdParser 解析Cmd命令
func (context *CQContext) CmdParser(message string, cmds ...string) (cmd string, params []string) {
	msg := strings.TrimSpace(message)
	split := strings.Split(msg, " ")
	tcmd := split[0]
	if tcmd == "" {
		return "", nil
	}
	for _, v := range cmds {
		if len(v) == len(tcmd) && v != tcmd || len(v) > len(tcmd) {
			continue
		}
		if v == tcmd[:len(v)] {
			if len(tcmd[len(v):]) == 0 {
				continue
			}
			cmd = tcmd[len(v):]
			if len(split) == 1 {
				params = nil
			} else {
				params = split[1:]
			}
			break
		}
	}
	return
}

// JSON quick response
func (context *CQContext) JSON(code int, body interface{}) error {
	if context.ws {
		return nil
	}
	if context.quick {
		return nil
	}
	context.Context.Writer.WriteHeader(code)
	if body == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}
	context.Context.Writer.Write(buf.Bytes())
	context.quick = true
	return nil
}

// User is the sender of cqevent
type User struct {
	ID       int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
	Area     string `json:"area,omitempty"`

	Card                string `json:"card,omitempty"`
	CardChangeable      bool   `json:"card_changeable,omitempty"`
	Title               string `json:"title,omitempty"`
	TitleExpireTimeUnix int64  `json:"title_expire_time,omitempty"`
	Level               string `json:"level,omitempty"`
	Role                string `json:"role,omitempty"`
	Unfriendly          bool   `json:"unfriendly,omitempty"`
	JoinTimeUnix        int64  `json:"join_time,omitempty"`
	LastSentTimeUnix    int64  `json:"last_sent_time,omitempty"`
	AnonymousID         int64  `json:"anonymous_id,omitempty" anonymous:"id"`
	AnonymousName       string `json:"anonymous_name,omitempty" anonymous:"name"`
	AnonymousFlag       string `json:"anonymous_flag,omitempty" anonymous:"flag"`
}

// File is the cqhttp event file
type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	BusID int64  `json:"busid"`
}
