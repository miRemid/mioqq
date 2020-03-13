package mioqq

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
)

const (
	// PrivateMessage 发送私人消息
	PrivateMessage = "send_private_msg"
	// PrivateMessageAsync 发送私人消息(异步)
	PrivateMessageAsync = "send_private_msg_async"
	// GroupMessage 发送群消息
	GroupMessage = "send_group_msg"
	// GroupMessageAsync 发送群消息(异步)
	GroupMessageAsync = "send_group_msg_async"
	// DiscussMessage 发送讨论组消息
	DiscussMessage = "send_discuss_msg"
	// DiscussMessageAsync 发送讨论组消息(异步)
	DiscussMessageAsync = "send_discuss_msg_async"
	// SendMessage 发送消息
	SendMessage = "send_msg"
	// SendMessageAsync 发送消息(异步)
	SendMessageAsync = "send_msg_async"
	// DeleteMessage 撤回消息
	DeleteMessage = "delete_msg"
	// SendLike 发送好友赞
	SendLike = "send_like"
	// SetGroupKick 踢人
	SetGroupKick = "set_group_kick"
	// SetGroupBan 禁言
	SetGroupBan = "set_group_ban"
	// SetGroupAnonymousBan 匿名禁言
	SetGroupAnonymousBan = "set_group_anonymous_ban"
	// SetGroupWholeBan 全体禁言
	SetGroupWholeBan = "set_group_whole_ban"
	// SetGroupAdmin 设置管理
	SetGroupAdmin = "set_group_admin"
	// SetGroupAnonymous 设置群匿名
	SetGroupAnonymous = "set_group_anonymous"
	// SetGroupCard 修改群名片
	SetGroupCard = "set_group_card"
	// SetGroupLeave 退群
	SetGroupLeave = "set_group_leave"
	// SetGroupSpecialTitle 设置群头衔
	SetGroupSpecialTitle = "set_group_special_title"
	// SetDiscussLeave 设置讨论组
	SetDiscussLeave = "set_discuss_leave"
	// SetFriendAddRequest 处理好友邀请
	SetFriendAddRequest = "set_friend_add_request"
	// SetGroupAddRequest 处理群邀请
	SetGroupAddRequest = "set_group_add_request"
	// GetLoginInfo 获取登录账号信息
	GetLoginInfo = "get_login_info"
	// GetStrangerInfo 获取陌生人信息
	GetStrangerInfo = "get_stranger_info"
	// GetFriendList 获取好友列表
	GetFriendList = "get_friend_list"
	// GetGroupList 获取群组列表
	GetGroupList = "get_group_list"
	// GetGroupInfo 获取群信息
	GetGroupInfo = "get_group_info"
	// getGroupMemberInfo 获取群成员信息
	getGroupMemberInfo = "get_group_member_info"
	// GetGroupMemberList 获取群成员列表
	GetGroupMemberList = "get_group_member_list"
	// GetCookies 获取Cookies
	GetCookies = "get_cookies"
	// GetCsrfToken 获取CsrfToken
	GetCsrfToken = "get_csrf_token"
	// GetCredentials 获取QQ凭证
	GetCredentials = "get_credentials"
	// GetRecord 获取语音
	GetRecord = "get_record"
	// GetImage 获取图片
	GetImage = "get_image"
	// CanSendImage 检查是否可以发送图片
	CanSendImage = "can_send_image"
	// CanSendRecord 检查是否可以发送语音
	CanSendRecord = "can_send_record"
	// GetStatus 获取插件状态
	GetStatus = "get_status"
	// GetVersionInfo 获取插件和COOLQ版本信息
	GetVersionInfo = "get_version_info"
	// SetRestartPlugin 重启插件
	SetRestartPlugin = "set_restart_plugin"
	// CleanDataDir 清理数据文件夹
	CleanDataDir = "clean_data_dir"
	// CleanPluginLog 清理日志文件夹
	CleanPluginLog = "clean_plugin_log"
)

// Send a message
func (m *API) Send(content *Message) (res CQAPIResponse, err error) {
	data, err := content.Parse()
	if err != nil {
		return CQAPIResponse{}, err
	}
	buff := bytes.NewBuffer(data)
	return m.Request(content.Type, buff)
}

// SendPrivateMsg will send a message to private user
func (m *API) SendPrivateMsg(to int64, content *Message, async bool) (res CQAPIResponse, err error) {
	content.ID = to
	if async {
		content.Type = PrivateMessageAsync
	} else {
		content.Type = PrivateMessage
	}
	return m.Send(content)
}

// SendGroupMsg will send a message to group
func (m *API) SendGroupMsg(to int64, content *Message, async bool) (res CQAPIResponse, err error) {
	content.ID = to
	if async {
		content.Type = GroupMessageAsync
	} else {
		content.Type = GroupMessage
	}
	return m.Send(content)
}

// SendDiscussMsg will send a message to discuss
func (m *API) SendDiscussMsg(to int64, content *Message, async bool) (res CQAPIResponse, err error) {
	content.ID = to
	if async {
		content.Type = DiscussMessageAsync
	} else {
		content.Type = DiscussMessage
	}
	return m.Send(content)
}

// SendMessage will send a message by cqhttp
func (m *API) SendMessage(to int64, content *Message, messageType string, async bool) (res CQAPIResponse, err error) {
	content.ID = to
	if async {
		content.Type = SendMessageAsync
	} else {
		content.Type = SendMessage
	}
	content.MessageType = messageType
	return m.Send(content)
}

// APIRequest send a normal request
func (m *API) APIRequest(api string, values *url.Values) (res CQAPIResponse, err error) {
	if values != nil {
		return m.Request(api, strings.NewReader(values.Encode()))
	}
	return m.Request(api, nil)
}

// DeleteMessage will delete id message
func (m *API) DeleteMessage(id int32) {
	values := url.Values{}
	values.Set("message_id", strconv.Itoa(int(id)))
	m.APIRequest(DeleteMessage, &values)
}

// SendLike 发送好友赞
func (m *API) SendLike(userid int64, times int) {
	values := url.Values{}
	values.Set("user_id", strconv.Itoa(int(userid)))
	values.Set("times", strconv.Itoa(times))
	m.APIRequest(SendLike, &values)
}

// SetGroupKick 群组踢人
func (m *API) SetGroupKick(groupid, userid int64, rejectAddRequest bool) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("user_id", strconv.Itoa(int(userid)))
	if rejectAddRequest {
		values.Set("reject_add_request", "true")
	}
	m.APIRequest(SetGroupKick, &values)
}

// SetGroupBan 单人禁言
func (m *API) SetGroupBan(groupid, userid, duration int64) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("user_id", strconv.Itoa(int(userid)))
	values.Set("duration", strconv.Itoa(int(duration)))
	m.APIRequest(SetGroupBan, &values)
}

// SetGroupAnonymousBan 群组匿名用户禁言
func (m *API) SetGroupAnonymousBan(groupid int64, anonymousFlag string, duration int64) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("duration", strconv.Itoa(int(duration)))
	values.Set("anonymous_flag", anonymousFlag)
	m.APIRequest(SetGroupAnonymousBan, &values)
}

// SetGroupWholeBan 全体禁言
func (m *API) SetGroupWholeBan(groupid int64, enable bool) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	if !enable {
		values.Set("enable", "false")
	}
	m.APIRequest(SetGroupWholeBan, &values)
}

// SetGroupAnonymous 群组匿名
func (m *API) SetGroupAnonymous(groupid int64, enable bool) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	if !enable {
		values.Set("enable", "false")
	}
	m.APIRequest(SetGroupAnonymous, &values)
}

// SetGroupCard 设置群名片
func (m *API) SetGroupCard(groupid, userid int64, card string) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("duration", strconv.Itoa(int(userid)))
	values.Set("card", card)
	m.APIRequest(SetGroupCard, &values)
}

// SetGroupLeave 退出群组
func (m *API) SetGroupLeave(groupid int64, isdismiss bool) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	if isdismiss {
		values.Set("is_dismiss", "true")
	}
	m.APIRequest(SetGroupLeave, &values)
}

// SetGroupSpecialTitle 设置群组专属头衔
func (m *API) SetGroupSpecialTitle(groupid, userid int64, specialTitle string, duration int64) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("duration", strconv.Itoa(int(userid)))
	values.Set("user_id", strconv.Itoa(int(userid)))
	values.Set("special_title", specialTitle)
	m.APIRequest(SetGroupSpecialTitle, &values)
}

// SetDiscussLeave 退出讨论组
func (m *API) SetDiscussLeave(discussid int64) {
	values := url.Values{}
	values.Set("discuss_id", strconv.Itoa(int(discussid)))
	m.APIRequest(SetDiscussLeave, &values)
}

// SetFriendAddRequest 处理加好友请求
func (m *API) SetFriendAddRequest(flag string, approve bool, remark string) {
	values := url.Values{}
	values.Set("flag", flag)
	if !approve {
		values.Set("approve", "false")
	}
	values.Set("remark", remark)
	m.APIRequest(SetFriendAddRequest, &values)
}

// SetGroupAddRequest 处理加群请求/邀请
func (m *API) SetGroupAddRequest(flag, subType string, approve bool, reason string) {
	values := url.Values{}
	values.Set("flag", flag)
	values.Set("sub_type", subType)
	if !approve {
		values.Set("approve", "false")
	}
	values.Set("reason", reason)
	m.APIRequest(SetGroupAddRequest, &values)
}

// GetLoginInfo 获取登录号信息
func (m *API) GetLoginInfo() (CQAPIResponse, error) {
	return m.APIRequest(GetLoginInfo, nil)
}

// GetStrangerInfo 获取陌生人信息
func (m *API) GetStrangerInfo(userid int64, noCache bool) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("user_id", strconv.Itoa(int(userid)))
	if noCache {
		values.Set("no_cache", "true")
	}
	return m.APIRequest(GetStrangerInfo, &values)
}

// GetFriendList 获取好友列表
func (m *API) GetFriendList() (CQAPIResponse, error) {
	return m.APIRequest(GetFriendList, nil)
}

// GetGroupList 获取群列表
func (m *API) GetGroupList() (CQAPIResponse, error) {
	return m.APIRequest(GetGroupList, nil)
}

// GetGroupInfo 获取群信息
func (m *API) GetGroupInfo(groupid int64, noCache bool) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	if noCache {
		values.Set("no_cache", "true")
	}
	return m.APIRequest(GetGroupInfo, &values)
}

// GetGroupMemberInfo 获取群成员信息
func (m *API) GetGroupMemberInfo(groupid, userid int64, noCache bool) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	values.Set("user_id", strconv.Itoa(int(userid)))
	if noCache {
		values.Set("no_cache", "true")
	}
	return m.APIRequest(getGroupMemberInfo, &values)
}

// GetGroupMemberList 获取群成员列表
func (m *API) GetGroupMemberList(groupid int64) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("group_id", strconv.Itoa(int(groupid)))
	return m.APIRequest(GetGroupMemberList, &values)
}

// GetCookies 获取Cookies
func (m *API) GetCookies(domain string) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("domain", domain)
	return m.APIRequest(GetCookies, &values)
}

// GetCsrfToken 获取CsrfToken
func (m *API) GetCsrfToken() (CQAPIResponse, error) {
	return m.APIRequest(GetCsrfToken, nil)
}

// GetCredentials 获取QQ相关接口凭证
func (m *API) GetCredentials(domain string) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("domain", domain)
	return m.APIRequest(GetCredentials, &values)
}

// GetRecord 获取语音
func (m *API) GetRecord(file, outFormat string, fullPath bool) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("file", file)
	values.Set("out_format", outFormat)
	if fullPath {
		values.Set("full_path", "true")
	}
	return m.APIRequest(GetRecord, &values)
}

// GetImage 获取图片
func (m *API) GetImage(file string) (CQAPIResponse, error) {
	values := url.Values{}
	values.Set("file", file)
	return m.APIRequest(GetImage, &values)
}

// CanSendImage 是否可以发送图片
func (m *API) CanSendImage() (CQAPIResponse, error) {
	return m.APIRequest(CanSendImage, nil)
}

// CanSendRecord 是否可以发送语音
func (m *API) CanSendRecord() (CQAPIResponse, error) {
	return m.APIRequest(CanSendRecord, nil)
}

// GetStatus 获取插件运行状态
func (m *API) GetStatus() (CQAPIResponse, error) {
	return m.APIRequest(GetStatus, nil)
}

// GetVersionInfo 获取COOLQ&CQHTTP版本信息
func (m *API) GetVersionInfo() (CQAPIResponse, error) {
	return m.APIRequest(GetVersionInfo, nil)
}

// SetRestartPlugin 重启插件
func (m *API) SetRestartPlugin() {
	m.APIRequest(SetRestartPlugin, nil)
}

// CleanDataDir 清理数据目录
func (m *API) CleanDataDir() {
	m.APIRequest(CleanDataDir, nil)
}

// CleanPluginLog 清理插件日志
func (m *API) CleanPluginLog() {
	m.APIRequest(CleanPluginLog, nil)
}
