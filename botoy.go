/* Package botoy 极简的OPQ开发框架
这是一个十分简单的框架，模块的方法很少，只封装了常用的几个功能，
比如发文字，发图。。。但在多数情况下够用了
*/
package botoy

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asmcos/requests"
	"github.com/goinggo/mapstructure"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var (
	ErrNotFunction      = errors.New("不是函数类型")
	ErrWrongArguments   = errors.New("参数有误，比如数目不对")
	ErrNotSupportedType = errors.New("不支持的类型")
	ErrInvalidBotID     = errors.New("无效的机器人QQ号")
)

type Botoy struct {
	// 机器人QQ号
	BotID int64
	// 连接地址 如: http://127.0.0.1:8888
	Address            string
	receiverCollection map[string][]reflect.Value
	mutex              sync.RWMutex
	stop               chan bool
	connected          bool
}

// NewBotoy 新建Botoy
//
// 如果bot为0，会监听所有机器人的消息，如果为一个确定的qq号，则只会接收该QQ号机器人所收到的消息
//
// 如果bot为0，将无法使用 Do 方法
func NewBotoy(bot int64, address string) *Botoy {
	return &Botoy{BotID: bot, Address: address, receiverCollection: make(map[string][]reflect.Value), mutex: sync.RWMutex{}, stop: make(chan bool, 1), connected: false}
}

// Start 启动
func (bot *Botoy) Start() (err error) {
	if bot.connected {
		return
	}
	log.Println("尝试连接中...")
	client, err := gosocketio.Dial(strings.ReplaceAll(strings.TrimSuffix(bot.Address, "/"), "http://", "ws://")+"/socket.io/?EIO=3&transport=websocket", transport.GetDefaultWebsocketTransport())
	if err != nil {
		return
	}
	// 连接
	client.On(gosocketio.OnConnection, func(_ *gosocketio.Channel) {
		bot.connected = true
		log.Println("已成功连接!")
		if receivers, ok := bot.receiverCollection[onConnectedPack]; ok {
			for _, receiver := range receivers {
				go receiver.Call([]reflect.Value{reflect.ValueOf(ConnectedPack{})})
			}
		}

	})
	// 断开连接
	client.On(gosocketio.OnDisconnection, func(_ *gosocketio.Channel) {
		bot.connected = false
		log.Println("连接已断开!")
		if receivers, ok := bot.receiverCollection[onDisConnectedPack]; ok {
			for _, receiver := range receivers {
				go receiver.Call([]reflect.Value{reflect.ValueOf(DisConnectedPack{})})
			}
		}
	restart:
		err := bot.Start()
		if err != nil {
			log.Println(err)
			log.Println("五秒后重连")
			time.Sleep(5 * time.Second)
			goto restart
		}
	})
	// 群消息
	client.On("OnGroupMsgs", func(_ *gosocketio.Channel, rawPack basePack) {
		if bot.BotID != 0 && rawPack.CurrentQQ != bot.BotID {
			return
		}
		log.Println("GroupMsg: ", rawPack)
		if receivers, ok := bot.receiverCollection[onGroupMsgsPack]; ok {
			packet := GroupMsgPack{}
			if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
				packet.CurrentQQ = rawPack.CurrentQQ
				for _, receiver := range receivers {
					go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
				}
			}
		}
	})
	// 好友消息
	client.On("OnFriendMsgs", func(_ *gosocketio.Channel, rawPack basePack) {
		if bot.BotID != 0 && rawPack.CurrentQQ != bot.BotID {
			return
		}
		log.Println("FriendMsg: ", rawPack)
		if receivers, ok := bot.receiverCollection[onFriendMsgsPack]; ok {
			packet := FriendMsgPack{}
			if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
				packet.CurrentQQ = rawPack.CurrentQQ
				for _, receiver := range receivers {
					go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
				}
			}
		}
	})
	// 事件
	client.On("OnEvents", func(_ *gosocketio.Channel, rawPack basePack) {
		if bot.BotID != 0 && rawPack.CurrentQQ != bot.BotID {
			return
		}
		log.Println("Events: ", rawPack)
		event, ok := rawPack.CurrentPacket.Data.(map[string]interface{})
		if !ok {
			return
		}
		eventName, ok := event["EventName"].(string)
		if !ok {
			return
		}
		switch eventName {
		case GroupJoinEventType:
			if receivers, ok := bot.receiverCollection[onGroupJoinPack]; ok {
				packet := GroupJoinPack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		case GroupExitEventType:
			if receivers, ok := bot.receiverCollection[onGroupExitPack]; ok {
				packet := GroupExitPack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		case GroupExitSuccEventType:
			if receivers, ok := bot.receiverCollection[onGroupExitSuccessPack]; ok {
				packet := GroupExitSuccessPack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		case GroupRevokeEventType:
			if receivers, ok := bot.receiverCollection[onGroupRevokePack]; ok {
				packet := GroupRevokePack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		case GroupShutEventType:
			if receivers, ok := bot.receiverCollection[onGroupShutPack]; ok {
				packet := GroupShutPack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		case FriendRevokeEventType:
			if receivers, ok := bot.receiverCollection[onFriendRevokePack]; ok {
				packet := FriendRevokePack{}
				if err = mapstructure.Decode(rawPack.CurrentPacket.Data, &packet); err == nil {
					packet.CurrentQQ = rawPack.CurrentQQ
					for _, receiver := range receivers {
						go receiver.Call([]reflect.Value{reflect.ValueOf(packet)})
					}
				}
			}
		default:
			log.Println("此类事件类型暂未处理")
		}
	})
	<-bot.stop
	return
}

// Stop 停止
func (bot *Botoy) Stop() {
	bot.stop <- true
}

// AddReceiver 添加对应消息的处理函数
/*
参数必须是函数，参数必须为模块内名字以Pack结尾的类型

可以添加多个函数，收到消息后各个函数异步执行

每个消息都可以添加多个函数进行处理

使用示例:
    bot.AddReceiver(
        func(pack botoy.FriendMsgPack) {
            fmt.Printf("收到一条好友消息， 发送人: %d, 内容: %s\n", pack.FromUin, pack.Content)
        },
        func(pack botoy.GroupMsgPack) {
            fmt.Printf("收到一条群消息， 发送人: %d, 群id: %d, 内容: %s\n", pack.FromUserID, pack.FromGroupID, pack.Content)
        },
        func(pack botoy.ConnectedPack) {
            fmt.Println("socketio已连接!")
        },
    )
    bot.AddReceiver(
        func(pack botoy.FriendRevokePack) {
            fmt.Printf("好友【%d】撤回了一条消息!\n", pack.Eventdata.Userid)
        },
    )
    bot.AddReceiver(
        func(pack botoy.GroupMsgPack) {
            fmt.Printf("收到一条群消息， 消息类型为: %s\n", pack.MsgType)
        },
    )
*/
func (bot *Botoy) AddReceiver(receiver ...interface{}) (err error) {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	for _, f := range receiver {
		fVal := reflect.ValueOf(f)
		if fVal.Kind() != reflect.Func {
			log.Println(ErrNotFunction)
			return ErrNotFunction
		}
		if fVal.Type().NumIn() != 1 {
			log.Println(ErrWrongArguments)
			return ErrWrongArguments
		}
		packName := fVal.Type().In(0).String()
		// 简单的做个判断
		if strings.HasPrefix(packName, "botoy.") && strings.HasSuffix(packName, "Pack") {
			if receivers, ok := bot.receiverCollection[packName]; ok {
				bot.receiverCollection[packName] = append(receivers, fVal)
			} else {
				bot.receiverCollection[packName] = []reflect.Value{fVal}
			}
		} else {
			log.Println(ErrNotSupportedType)
			return ErrNotSupportedType
		}
	}
	return
}

// Do 进行发送等操作
//
// 接收参数为模块内以Payload结尾的结构体，每个结构体内的字段无需全部指定
//
// 对于哪些是必须指定的字段，需要你自己十分了解机器人服务端api的使用
func (bot *Botoy) Do(payload payload) (resp *requests.Response, err error) {
	if bot.BotID == 0 {
		return nil, ErrInvalidBotID
	}
	method, apiPath, funcName, jsonStr := payload.process()
	// 处理请求地址
	address := strings.TrimSuffix(bot.Address, "/")
	apiPath = strings.TrimPrefix(apiPath, "/")
	api, _ := url.Parse(address + "/" + apiPath)
	params := url.Values{}
	params.Set("qq", strconv.FormatInt(bot.BotID, 10))
	params.Set("funcname", funcName)
	params.Set("timeout", "20")
	api.RawQuery = params.Encode()

	req := requests.Requests()
	req.SetTimeout(20)
	if method == http.MethodPost {
		resp, err = req.PostJson(api.String(), jsonStr)
		if err != nil {
			return
		}
		res := struct {
			Msg string `json:"Msg"`
			Ret int    `json:"Ret"`
		}{}
		if err1 := resp.Json(&res); err1 == nil {
			switch res.Ret {
			case 0:
			case 34:
				log.Println("未知错误，跟消息长度似乎无关，可以尝试分段重新发送 => ", res.Ret, res.Msg)
			case 110:
				log.Println("发送失败，你已被移出该群，请重新加群 => ", res.Ret, res.Msg)
			case 120:
				log.Println("机器人被禁言 => ", res.Ret, res.Msg)
			case 241:
				log.Println("消息发送频率过高，对同一个群或好友，建议发消息的最小间隔控制在1100ms以上 => ", res.Ret, res.Msg)
			case 299:
				log.Println("超过群发言频率限制 => ", res.Ret, res.Msg)
			default:
				log.Println("请求发送成功, 但处理失败 => ", res.Ret, res.Msg)
			}
		}
	}
	return
}
