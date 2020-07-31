package service

import (
	"container/list"
	"encoding/json"
	"github.com/beatrice950201/araneid/extend/model/inform"
	"github.com/gorilla/websocket"
)

type DefaultSocketService struct {
	Subscribe           chan DefaultSocketSubscriber // 新加入用户的chan。
	Unsubscribe         chan int                     // 离开用户的chan。
	Command             chan map[string]interface{}  // 需要处理的指令 比如爬虫指令等
	Message             chan inform.Message          // 需要处理的指令 比如爬虫指令等
	Subscribers         *list.List                   // 目前的用户列表
	collectService      DefaultCollectService        // 采集层逻辑
	dictionariesService DefaultDictionariesService   // 词典采集逻辑层
	movieService        DefaultMovieService          // 剧情采集逻辑层
}

/** socket常驻内存实例 **/
var socketInstanceObject *DefaultSocketService

/*链接存储器*/
type DefaultSocketSubscriber struct {
	User int
	Conn *websocket.Conn
}

/** 加入 **/
func (service *DefaultSocketService) Join(user int, ws *websocket.Conn) {
	service.Subscribe <- DefaultSocketSubscriber{User: user, Conn: ws}
}

/** 离开 **/
func (service *DefaultSocketService) Leave(user int) {
	service.Unsubscribe <- user
}

/** 发送指令 **/
func (service *DefaultSocketService) InstructHandle(instruct map[string]interface{}) {
	service.Command <- instruct
}

/** 发送通知 **/
func (service *DefaultSocketService) InformHandle(message inform.Message) {
	service.Message <- message
}

/** 常轮询监听动作 **/
func (service *DefaultSocketService) ServiceRoom() {
	for {
		select {
		case sub := <-service.Subscribe: // 用户加入
			service.userPush(service.Subscribers, sub)
		case msg := <-service.Message: // 处理通知
			service.messageHandle(msg)
		case instruct := <-service.Command: // 处理指令
			service.commandHandle(instruct)
		case unsub := <-service.Unsubscribe: // 用户离开
			service.userExit(service.Subscribers, unsub)
		}
	}
}

/** 处理指令 **/
func (service *DefaultSocketService) commandHandle(instruct map[string]interface{}) {
	switch instruct["command"].(string) {
	case "collect": // 采集器移交主要服务层处理
		go service.collectService.InstanceBegin(instruct)
	case "dict": // 词典采集器移交主要服务层处理
		go service.dictionariesService.CateInstanceBegin(instruct)
	case "movie": // 剧情采集器移交主要服务层处理
		go service.movieService.CateInstanceBegin(instruct)
	}
}

/** 处理通知 **/
func (service *DefaultSocketService) messageHandle(msg inform.Message) {
	for sub := service.Subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(DefaultSocketSubscriber).User == msg.Receiver {
			if ws := sub.Value.(DefaultSocketSubscriber).Conn; ws != nil {
				data, _ := json.Marshal(msg)
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					service.Unsubscribe <- sub.Value.(DefaultSocketSubscriber).User
				}
			}
			break
		}
	}
}

/** 用户是否已经存在 **/
func (service *DefaultSocketService) userPush(subscribers *list.List, sub DefaultSocketSubscriber) {
	var result bool
	for s := subscribers.Front(); s != nil; s = s.Next() {
		if s.Value.(DefaultSocketSubscriber).User == sub.User {
			result = true
		}
	}
	if !result {
		service.Subscribers.PushBack(sub)
	}
}

/**  用户离开 **/
func (service *DefaultSocketService) userExit(subscribers *list.List, unsub int) {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(DefaultSocketSubscriber).User == unsub {
			subscribers.Remove(sub)
			if ws := sub.Value.(DefaultSocketSubscriber).Conn; ws != nil {
				_ = ws.Close()
			}
			break
		}
	}
}

/** 获取到socket实例 **/
func SocketInstanceGet() *DefaultSocketService {
	if socketInstanceObject == nil {
		socketInstanceObject = &DefaultSocketService{
			Subscribe:   make(chan DefaultSocketSubscriber, 10),
			Unsubscribe: make(chan int, 10),
			Command:     make(chan map[string]interface{}, 10),
			Message:     make(chan inform.Message, 10),
			Subscribers: list.New(),
		}
	}
	return socketInstanceObject
}
