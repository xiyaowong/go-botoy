package botoy

type basePack struct {
	CurrentPacket struct {
		Data      interface{} `json:"Data"`
		WebConnID string      `json:"WebConnId"`
	} `json:"CurrentPacket"`
	CurrentQQ int64 `json:"CurrentQQ"`
}

// ConnectedPack socketio连接上
type ConnectedPack struct{}

const onConnectedPack = "botoy.ConnectedPack"

// DisConnectedPack socketio 断开连接
type DisConnectedPack struct{}

const onDisConnectedPack = "botoy.DisConnectedPack"

// GroupMsgPack 群消息
type GroupMsgPack struct {
	CurrentQQ     int64       `json:"CurrentQQ,omitempty"`
	Content       string      `json:"Content"`
	FromGroupID   int64       `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int         `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int         `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}

const onGroupMsgsPack = "botoy.GroupMsgPack"

// FriendMsgPack 好友消息
type FriendMsgPack struct {
	CurrentQQ  int64       `json:"CurrentQQ,omitempty"`
	Content    string      `json:"Content"`
	FromUin    int64       `json:"FromUin"`
	MsgSeq     int         `json:"MsgSeq"`
	MsgType    string      `json:"MsgType"`
	ToUin      int64       `json:"ToUin"`
	TempUin    int64       `json:"TempUin"`
	RedBaginfo interface{} `json:"RedBaginfo"`
}

const onFriendMsgsPack = "botoy.FriendMsgPack"

// GroupJoinPack 加群事件
type GroupJoinPack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	EventData struct {
		InviteUin int64  `json:"InviteUin"`
		UserID    int64  `json:"UserID"`
		UserName  string `json:"UserName"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

const onGroupJoinPack = "botoy.GroupJoinPack"

// GroupExitPack 退群事件
type GroupExitPack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	EventData struct {
		UserID int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

const onGroupExitPack = "botoy.GroupExitPack"

// GroupExitSuccessPack 退群事件
type GroupExitSuccessPack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	EventData struct {
		GroupID int64 `json:"GroupID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

const onGroupExitSuccessPack = "botoy.GroupExitSuccessPack"

// GroupRevokePack 群撤回消息
type GroupRevokePack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	EventData struct {
		AdminUserID int   `json:"AdminUserID"`
		GroupID     int64 `json:"GroupID"`
		MsgRandom   int64 `json:"MsgRandom"`
		MsgSeq      int   `json:"MsgSeq"`
		UserID      int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

const onGroupRevokePack = "botoy.GroupRevokePack"

// GroupShutPack 群禁言
type GroupShutPack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	EventData struct {
		GroupID  int64 `json:"GroupID"`
		ShutTime int   `json:"ShutTime"`
		UserID   int64 `json:"UserID"`
	} `json:"EventData"`
	EventMsg struct {
		FromUin    int64       `json:"FromUin"`
		ToUin      int64       `json:"ToUin"`
		MsgType    string      `json:"MsgType"`
		MsgSeq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		RedBaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
}

const onGroupShutPack = "botoy.GroupShutPack"

// FriendRevokePack 好友撤回事件
type FriendRevokePack struct {
	CurrentQQ int64 `json:"CurrentQQ,omitempty"`
	Eventdata struct {
		Msgseq int `json:"MsgSeq"`
		Userid int `json:"UserID"`
	} `json:"EventData"`
	Eventmsg struct {
		Fromuin    int         `json:"FromUin"`
		Touin      int         `json:"ToUin"`
		Msgtype    string      `json:"MsgType"`
		Msgseq     int         `json:"MsgSeq"`
		Content    string      `json:"Content"`
		Redbaginfo interface{} `json:"RedBaginfo"`
	} `json:"EventMsg"`
	Eventname string `json:"EventName"`
}

const onFriendRevokePack = "botoy.FriendRevokePack"
