package cqcode

import (
	"fmt"
	"strings"
)

// CQCode create a string cqcode
func CQCode(name string, params ...interface{}) string {
	cqcode := fmt.Sprintf("[CQ:%s,", name)
	var i int = 0
	for {
		if i == len(params) {
			return strings.TrimRight(cqcode, ",") + "]"
		}
		cqcode += fmt.Sprintf("%v=%v,", params[i], params[i+1])
		i += 2
	}
}

// Face 返回自带表情cqcode
func Face(id int) string {
	return fmt.Sprintf("[CQ:face,id=%d]", id)
}

// BFace 返回原创表情cqcode
// 原创表情存放在酷Q目录的data\bface\下
func BFace(id int) string {
	return fmt.Sprintf("[CQ:bface,id=%d]", id)
}

// SFace 小表情
func SFace(id int) string {
	return fmt.Sprintf("[CQ:sface,id=%d]", id)
}

// Emoji 返回emoji的cqcode
func Emoji(unicode int) string {
	return fmt.Sprintf("[CQ:emoji,id=%d]", unicode)
}

// Image 图片cq码
// uri支持本地图片或url链接，同时支持base64编码，需base64://前缀
// cache为是否使用缓存
func Image(uri string, cache bool) string {
	if cache {
		return fmt.Sprintf("[CQ:image,file=%s]", uri)
	}
	return fmt.Sprintf("[CQ:image,cache=0,file=%s]", uri)
}

// Recode 语音cq码
// uri支持本地图片或url链接，同时支持base64编码，需base64://前缀
// cache为是否使用缓存
func Recode(uri string, cache bool) string {
	if cache {
		return fmt.Sprintf("[CQ:recode,file=%s]", uri)
	}
	return fmt.Sprintf("[CQ:recode,cache=0,file=%s]", uri)
}

// At @某人
func At(userid int) string {
	return fmt.Sprintf("[CQ:at,qq=%d]", userid)
}

// AtAll @全体成员
func AtAll() string {
	return fmt.Sprintf("[CQ:at,qq=all]")
}

// RPS 猜拳魔法表情
func RPS(id int) string {
	return fmt.Sprintf("[CQ:rps,type=%d]", id)
}

// Dice 掷骰子表情
func Dice(id int) string {
	return fmt.Sprintf("[CQ:dice,type=%d]", id)
}

// Shake 戳一戳，仅好友消息
func Shake() string {
	return fmt.Sprintf("[CQ:shake]")
}

// Anonymous 匿名发送消息
// 需要加在消息开头
func Anonymous(ignore bool) string {
	if ignore {
		return fmt.Sprintf("[CQ:anonymous,ignore=true]")
	}
	return fmt.Sprintf("[CQ:anonymous]")
}

// Location 位置消息
func Location(lat, lon float64, title, content string) string {
	return fmt.Sprintf("[CQ:location,lat=%v,lon=%v,title=%s,content=%s]", lat, lon, title, content)
}

// Sign 签到
func Sign(location, title, imageurl string) string {
	return fmt.Sprintf("[CQ:sign,location=%s,title=%s,image=%s]", location, title, imageurl)
}

// ClientMusic 发送平台音乐
// client为音乐平台，目前支持qq、163
// id为对应音乐平台的数字音乐id
func ClientMusic(client string, id int, style string) string {
	return CQCode(
		"music",
		"type", client,
		"id", id,
		"style", style,
	)
}

// Music 发送自定义音乐分享
// url: 分享链接，点击分享后进入的音乐页面
// audio: 音频链接
// title: 标题
// content: 简介
// image: 封面链接，为空显示默认图片
func Music(url, audio, title, content, image string) string {
	return CQCode(
		"music",
		"type", "custom",
		"url", url,
		"audio", audio,
		"title", title,
		"content", content,
		"image", image,
	)
}

// Share 发送链接分享
// url: 分享的链接
// title: 标题
// content: 简介
// image: 图片链接，为空显示默认图片
func Share(url, title, content, image string) string {
	return CQCode(
		"share",
		"url", url,
		"title", title,
		"content", content,
		"image", image,
	)
}
