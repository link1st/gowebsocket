/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 14:38
 */

package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gowebsocket/common"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"sync"
	"time"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code int, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

func init() {
	Register("login", login)
	Register("heartbeat", heartbeat)
}

// 处理数据
func ProcessData(client *Client, message []byte) {

	fmt.Println("处理数据", client.Addr, string(message))

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()

	request := &models.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("处理数据 json Unmarshal", err)
		client.SendMsg([]byte("数据不合法"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		client.SendMsg([]byte("处理数据失败"))

		return
	}

	seq := request.Seq
	cmd := request.Cmd

	// 返回数据
	head := &models.Head{
		Seq: seq,
		Cmd: cmd,
	}

	var (
		code int
		msg  string
		data interface{}
	)

	// request
	fmt.Println("acc_request", cmd, client.Addr)

	// 采用 map 注册的方式
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		fmt.Println("处理数据 路由不存在", client.Addr, "cmd", cmd)
	}

	msg = common.GetErrorMessage(code, msg)

	response := models.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	head.Response = response

	headByte, err := json.Marshal(head)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)

		return
	}

	client.Send <- headByte

	fmt.Println("acc_response send", client.Addr, client.AppId, client.UserId, "cmd", cmd, "code", code)

	return
}

// 用户登录
func login(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", request.ServiceToken)

	if request.UserId == "" || len(request.UserId) >= 10 {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)

		return
	}

	client.Login(request.AppId, request.UserId, currentTime)

	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)

		return
	}

	// 用户登录
	login := &Login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Client: client,
	}
	clientManager.Login <- login

	fmt.Println("用户登录 成功", seq, client.Addr, request.UserId)

	return
}

// 心跳接口
func heartbeat(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 心跳接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.AppId, client.UserId, seq)
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GerUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			fmt.Println("心跳接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = common.ServerError
			fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)

			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	return
}
