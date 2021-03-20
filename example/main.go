package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/xiyaowong/go-botoy"
)

var (
	botID      int
	botAddress string
	bot        *botoy.Botoy
)

func init() {
	botID, err := strconv.Atoi(os.Getenv("BotID"))
	if err != nil {
		fmt.Print("输入机器人QQ号：")
		fmt.Scanf("%d", &botID)
	}
	botAddress = os.Getenv("BotAddress")
	if botAddress == "" {
		fmt.Print("输入机器人连接地址：")
		fmt.Scanf("%s", &botAddress)
	}
	fmt.Println(botID, botAddress)
	bot = botoy.NewBotoy(int64(botID), botAddress)
}

// 这个函数用来处理进群事件
func someoneJoin(pack botoy.GroupJoinPack) {
	bot.Do(&botoy.GroupTextPayload{GroupID: pack.EventMsg.FromUin, Text: fmt.Sprintf("欢迎 <%s> 加入本群!", pack.EventData.UserName)})
}

func main() {
	// 连接成功后给手机发送提示
	bot.AddReceiver(func(_ botoy.ConnectedPack) {
		bot.Do(&botoy.PhoneTextPayload{Text: "socketio已连接!"})
	})

	bot.AddReceiver(
		// 复读好友发给机器人的文字消息
		func(pack botoy.FriendMsgPack) {
			if pack.FromUin == pack.CurrentQQ {
				return
			}
			if pack.MsgType == botoy.TextMsgType {
				bot.Do(&botoy.FriendTextPayload{ToUserID: pack.FromUin, Text: pack.Content})
			}
		},
		// 收到群文字消息为 go 时，回复该消息，回复内容为 gogogo
		func(pack botoy.GroupMsgPack) {
			if pack.FromUserID == pack.CurrentQQ {
				return
			}
			if pack.MsgType == botoy.TextMsgType && pack.Content == "go" {
				bot.Do(&botoy.GroupReplyPayload{
					GroupID:       pack.FromGroupID,
					Content:       "gogogo",
					RawMsgContent: pack.Content,
					RawMsgSeq:     pack.MsgSeq,
					RawMsgTime:    pack.MsgTime,
					RawMsgUserID:  pack.FromUserID,
				})
			}
		},
	)

	// 好友撤回一条消息时，立即问他撤回了什么见不得人的东西
	bot.AddReceiver(
		func(pack botoy.FriendRevokePack) {
			bot.Do(&botoy.FriendTextPayload{ToUserID: int64(pack.Eventmsg.Fromuin), Text: "撤回了什么见不得人的东西??"})
		},
		// 新成员加群事件
		someoneJoin,
	)

	// 收到群消息为 壁纸 时 发送一张图片
	bot.AddReceiver(func(pack botoy.GroupMsgPack) {
		if pack.Content == "壁纸" {
			bot.Do(&botoy.GroupPicPayload{GroupID: pack.FromGroupID, PicURL: "http://api.btstu.cn/sjbz/?lx=dongman"})
		}
	})

	// 收到群消息为 base64 时 通过base64发送一张图片
	bot.AddReceiver(func(pack botoy.GroupMsgPack) {
		if pack.Content == "base64" {
			picBytes, err := ioutil.ReadFile("./base64.jpg")
			if err != nil {
				log.Println(err)
				return
			}
			picBase64 := base64.StdEncoding.EncodeToString(picBytes)
			bot.Do(&botoy.GroupPicPayload{GroupID: pack.FromGroupID, PicBase64: picBase64})
		}
	})

	/////////////////////////////////////
	if err := bot.Start(); err != nil {
		panic(err)
	}
}
