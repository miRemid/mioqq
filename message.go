package mioqq

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/miRemid/mioqq/cqcode"
)

const (
	// StringContent is the flag
	StringContent = "string"
	// ArrayContent is the flag
	ArrayContent = "block"
)

func block(t string, argv ...string) CQParams {
	var res = make(CQParams)
	res["type"] = t
	var data = make(CQParams)
	if len(argv) < 2 || len(argv)%2 != 0 {
		return res
	}
	i := 0
	for {
		if i == len(argv) {
			break
		}
		data[argv[i]] = argv[i+1]
		i += 2
	}
	res["data"] = data
	return res
}

// Message is the qq content struct
type Message struct {
	Async bool
	ID    int64

	MessageType   string
	MessageFormat string // 消息格式，string array
	Type          string // 消息类型， private， group, discuss

	stringMessage string
	arrayMessage  []map[string]interface{}
}

// NewMessage return a mio message
func (m *API) NewMessage(to int64, flag string, MessageFormat string) *Message {
	return &Message{
		Type:          flag,
		ID:            to,
		MessageFormat: MessageFormat,
		arrayMessage:  make([]map[string]interface{}, 0),
	}
}

// Content is the struct of mio's send message
type Content struct {
	Auto        bool        `json:"auto_escape"`
	Message     interface{} `json:"message"`
	UserID      int64       `json:"user_id,omitempty"`
	GroupID     int64       `json:"group_id,omitempty"`
	DiscussID   int64       `json:"discuss_id,omitempty"`
	MessageType string      `json:"message_type,omitempty"`
}

// Parse will return the []byte data
func (msg *Message) Parse() ([]byte, error) {
	if msg.Type == "" {
		msg.Type = PrivateMessage
	}
	var content Content
	switch msg.Type {
	case PrivateMessage:
		fallthrough
	case PrivateMessageAsync:
		content.UserID = msg.ID
		break
	case GroupMessage:
		fallthrough
	case GroupMessageAsync:
		content.GroupID = msg.ID
		break
	case DiscussMessage:
		fallthrough
	case DiscussMessageAsync:
		content.DiscussID = msg.ID
		break
	case SendMessage:
		switch msg.MessageType {
		case "private":
			content.UserID = msg.ID
			break
		case "group":
			content.GroupID = msg.ID
			break
		case "discuss":
			content.DiscussID = msg.ID
			break
		default:
			msg.MessageType = "private"
			content.UserID = msg.ID
			break
		}
		break
	}
	if msg.MessageFormat == ArrayContent {
		content.Message = msg.arrayMessage
	} else {
		content.Message = msg.stringMessage
	}
	if msg.MessageType != "" {
		content.MessageType = msg.MessageType
	}
	return json.Marshal(content)
}

// NewLine will add a \n
func (msg *Message) NewLine() *Message {
	return msg.Text("\n")
}

// Text will add a content into the message struct
func (msg *Message) Text(content string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("text", "text", content))
	} else {
		msg.stringMessage += content
	}
	return msg
}

// Emoji will add a Emoji cqcode into the message sturct
func (msg *Message) Emoji(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("emoji", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.Emoji(id)
	}
	return msg
}

// Face will add a Face cqcode into the message sturct
func (msg *Message) Face(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("face", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.Face(id)
	}
	return msg
}

// BFace will add a BFace cqcode into the message sturct
func (msg *Message) BFace(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("bface", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.BFace(id)
	}
	return msg
}

// SFace will add a SFace cqcode into the message sturct
func (msg *Message) SFace(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("sface", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.SFace(id)
	}
	return msg
}

// Image will add a Image cqcode into the message sturct
func (msg *Message) Image(filename string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("image", "file", filename))
	} else {
		msg.stringMessage += fmt.Sprintf("[CQ=image,file=%s]", filename)
	}
	return msg
}

// Recode will add a Recode cqcode into the message sturct
func (msg *Message) Recode(filename string, magic bool) *Message {
	if msg.MessageFormat == ArrayContent {
		if magic {
			msg.arrayMessage = append(msg.arrayMessage, block("recode", "file", filename, "magic", "true"))
		} else {
			msg.arrayMessage = append(msg.arrayMessage, block("recode", "file", filename))
		}
	} else {
		if magic {
			msg.stringMessage += fmt.Sprintf("[CQ=recode,file=%s,magic=true]", filename)
		} else {
			msg.stringMessage += fmt.Sprintf("[CQ=recode,file=%s]", filename)
		}
	}
	return msg
}

// RCCode is the cqhttp code plus plus plus~
func (msg *Message) RCCode(t, url string, cache int) *Message {
	if msg.MessageFormat == ArrayContent {
		if cache == 0 {
			msg.arrayMessage = append(msg.arrayMessage, block(t, "file", url, "cache", "0"))
		} else {
			msg.arrayMessage = append(msg.arrayMessage, block(t, "file", url))
		}
	} else {
		msg.stringMessage += fmt.Sprintf("[CQ=%s,file=%s,cache=%d]", t, url, cache)
	}
	return msg
}

// At somebody
func (msg *Message) At(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("at", "qq", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.At(id)
	}
	return msg
}

// AtAll at all of people
func (msg *Message) AtAll() *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("at", "qq", "all"))
	} else {
		msg.stringMessage += cqcode.AtAll()
	}
	return msg
}

// RPS will add a RPS cqcode into the message sturct
func (msg *Message) RPS(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("rps", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.RPS(id)
	}
	return msg
}

// Dice will add a Dice cqcode into the message sturct
func (msg *Message) Dice(id int) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("dice", "id", strconv.Itoa(id)))
	} else {
		msg.stringMessage += cqcode.Dice(id)
	}
	return msg
}

// Shake will add a Shake cqcode into the message sturct
func (msg *Message) Shake() *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block("shake"))
	} else {
		msg.stringMessage += cqcode.Shake()
	}
	return msg
}

// Anonymous will add a anonymous cqcode into the message sturct
func (msg *Message) Anonymous(ignore bool) *Message {
	if msg.MessageFormat == ArrayContent {
		if ignore {
			msg.arrayMessage = append(msg.arrayMessage, block("anonymous", "ignore", "true"))
		} else {
			msg.arrayMessage = append(msg.arrayMessage, block("anonymous"))
		}
	} else {
		msg.stringMessage += cqcode.Anonymous(ignore)
	}
	return msg
}

// Location will add a Location cqcode into the message sturct
func (msg *Message) Location(lat, lon float64, title, content string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block(
			"Location",
			"lat", strconv.FormatFloat(lat, 'f', 6, 64),
			"lon", strconv.FormatFloat(lon, 'f', 6, 64),
			"title", title,
			"content", content,
		))
	} else {
		msg.stringMessage += cqcode.Location(lat, lon, title, content)
	}
	return msg
}

// Sign will add a Sign cqcode into the message sturct
func (msg *Message) Sign(location, title, imageurl string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block(
			"sign",
			"location", location,
			"title", title,
			"image", imageurl,
		))
	} else {
		msg.stringMessage += cqcode.Sign(location, title, imageurl)
	}
	return msg
}

// ClientMusic will add a ClientMusic cqcode into the message sturct
func (msg *Message) ClientMusic(t string, id int, style string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block(
			"music",
			"type", t,
			"id", strconv.Itoa(id),
			"style", style,
		))
	} else {
		msg.stringMessage += cqcode.ClientMusic(t, id, style)
	}
	return msg
}

// Music will add a Music cqcode into the message sturct
func (msg *Message) Music(url, audio, title, content, image string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block(
			"music",
			"type", "custom",
			"url", url,
			"audio", audio,
			"title", title,
			"content", content,
			"image", image,
		))
	} else {
		msg.stringMessage += cqcode.Music(url, audio, title, content, image)
	}
	return msg
}

// Share will add a Share cqcode into the message sturct
func (msg *Message) Share(url, title, content, image string) *Message {
	if msg.MessageFormat == ArrayContent {
		msg.arrayMessage = append(msg.arrayMessage, block(
			"music",
			"url", url,
			"title", title,
			"content", content,
			"image", image,
		))
	} else {
		msg.stringMessage += cqcode.Share(url, title, content, image)
	}
	return msg
}

// CQMessage is the struct from cqhttp's response
type CQMessage struct {
	ID int32 `json:"message_id"`
}
