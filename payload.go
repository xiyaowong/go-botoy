package botoy

import (
	"encoding/json"
	"net/http"
)

type payload interface {
	// 返回 request_method, api_path, funcname, json_string
	process() (string, string, string, string)
}

// FriendTextPayload 好友文字消息
type FriendTextPayload struct {
	ToUserID int64
	Text     string
}

func (p *FriendTextPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 1
	jsonData["SendMsgType"] = "TextMsg"
	jsonData["ToUserUid"] = p.ToUserID
	jsonData["Content"] = p.Text
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}

// FriendPicPayload 好友图片消息
type FriendPicPayload struct {
	ToUser    int64
	Content   string
	PicURL    string
	PicBase64 string
	FileMd5   string
	Flash     bool
}

func (p *FriendPicPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["sendToType"] = 1
	jsonData["sendMsgType"] = "PicMsg"
	jsonData["toUser"] = p.ToUser
	jsonData["content"] = p.Content
	jsonData["picUrl"] = p.PicURL
	jsonData["picBase64Buf"] = p.PicBase64
	jsonData["fileMd5"] = p.FileMd5
	jsonData["flashPic"] = p.Flash
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsg", string(jsonBytes)
}

// FriendVoicePayload 好友语音消息
type FriendVoicePayload struct {
	ToUser      int64
	VoiceURL    string
	VoiceBase64 string
}

func (p *FriendVoicePayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["sendToType"] = 1
	jsonData["sendMsgType"] = "VoiceMsg"
	jsonData["toUser"] = p.ToUser
	jsonData["voiceUrl"] = p.VoiceURL
	jsonData["voiceBase64Buf"] = p.VoiceBase64
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsg", string(jsonBytes)
}

// GroupTextPayload 群文字消息
type GroupTextPayload struct {
	GroupID int64
	Text    string
}

func (p *GroupTextPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 2
	jsonData["SendMsgType"] = "TextMsg"
	jsonData["ToUserUid"] = p.GroupID
	jsonData["Content"] = p.Text
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}

// GroupPicPayload 群图片消息
type GroupPicPayload struct {
	GroupID   int64
	Content   string
	PicURL    string
	PicBase64 string
	PicMD5s   []string
	PicMD5    string
	Flash     bool
}

func (p *GroupPicPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	// v1
	jsonData["toUser"] = p.GroupID
	jsonData["sendToType"] = 2
	jsonData["sendMsgType"] = "PicMsg"
	jsonData["content"] = p.Content
	jsonData["picUrl"] = p.PicURL
	jsonData["picBase64Buf"] = p.PicBase64
	jsonData["fileMd5"] = p.PicMD5
	jsonData["flashPic"] = p.Flash
	// v2
	jsonData["ToUserUid"] = p.GroupID
	jsonData["SendToType"] = 2
	jsonData["SendMsgType"] = "PicMsg"
	jsonData["PicMd5s"] = p.PicMD5s
	jsonBytes, _ := json.Marshal(jsonData)

	var funcname string
	if len(p.PicMD5s) > 0 {
		funcname = "SendMsgV2"
	} else {
		funcname = "SendMsg"
	}
	return http.MethodPost, "v1/LuaApiCaller", funcname, string(jsonBytes)
}

// GroupVoicePayload 群语音消息
type GroupVoicePayload struct {
	GroupID     int64
	VoiceURL    string
	VoiceBase64 string
}

func (p *GroupVoicePayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["sendToType"] = 2
	jsonData["sendMsgType"] = "VoiceMsg"
	jsonData["toUser"] = p.GroupID
	jsonData["voiceUrl"] = p.VoiceURL
	jsonData["voiceBase64Buf"] = p.VoiceBase64
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsg", string(jsonBytes)
}

// GroupJSONPayload 群json消息
type GroupJSONPayload struct {
	GroupID int64
	JSONStr string
}

func (p *GroupJSONPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 2
	jsonData["SendMsgType"] = "JsonMsg"
	jsonData["ToUserUid"] = p.GroupID
	jsonData["Content"] = p.JSONStr
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}

// GroupXMLPayload 群XML消息
type GroupXMLPayload struct {
	GroupID int64
	XMLStr  string
}

func (p *GroupXMLPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 2
	jsonData["SendMsgType"] = "XmlMsg"
	jsonData["ToUserUid"] = p.GroupID
	jsonData["Content"] = p.XMLStr
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}

// GroupReplyPayload 群回复消息
type GroupReplyPayload struct {
	GroupID       int64
	Content       string
	RawMsgSeq     int
	RawMsgTime    int
	RawMsgUserID  int64
	RawMsgContent string
}

func (p *GroupReplyPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["sendToType"] = 2
	jsonData["sendMsgType"] = "ReplayMsg"
	jsonData["toUser"] = p.GroupID
	jsonData["content"] = p.Content
	jsonData["replayInfo"] = map[string]interface{}{
		"MsgSeq":     p.RawMsgSeq,
		"MsgTime":    p.RawMsgTime,
		"UserID":     p.RawMsgUserID,
		"RawContent": p.RawMsgContent,
	}
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsg", string(jsonBytes)
}

// PrivateTextPayload 私聊文字消息
type PrivateTextPayload struct {
	ToUserID int64
	GroupID  int64
	Text     string
}

func (p *PrivateTextPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 3
	jsonData["SendMsgType"] = "TextMsg"
	jsonData["GroupID"] = p.GroupID
	jsonData["ToUserUid"] = p.ToUserID
	jsonData["Content"] = p.Text
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}

// PrivatePicPayload 私聊图片消息
type PrivatePicPayload struct {
	ToUserID  int64
	GroupID   int64
	PicURL    string
	PicBase64 string
	FileMd5   string
	Content   string
	Flash     bool
}

func (p *PrivatePicPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["sendToType"] = 3
	jsonData["sendMsgType"] = "PicMsg"
	jsonData["groupid"] = p.GroupID
	jsonData["toUser"] = p.ToUserID
	jsonData["content"] = p.Content
	jsonData["picUrl"] = p.PicURL
	jsonData["picBase64Buf"] = p.PicBase64
	jsonData["fileMd5"] = p.FileMd5
	jsonData["flashPic"] = p.Flash
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsg", string(jsonBytes)
}

// PhoneTextPayload 给手机发文字消息
type PhoneTextPayload struct {
	Text string
}

func (p *PhoneTextPayload) process() (string, string, string, string) {
	jsonData := make(map[string]interface{})
	jsonData["SendToType"] = 2
	jsonData["SendMsgType"] = "PhoneMsg"
	jsonData["Content"] = p.Text
	jsonBytes, _ := json.Marshal(jsonData)
	return http.MethodPost, "v1/LuaApiCaller", "SendMsgV2", string(jsonBytes)
}
