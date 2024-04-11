// Package websocket 处理
package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"

	"github.com/redis/go-redis/v9"
)

// PingController ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)
	data = "pong"
	return
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)
		return
	}
	fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", request.ServiceToken)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID
	if request.UserID == "" || len(request.UserID) >= 20 {
		code = common.UnauthorizedUserID
		fmt.Println("用户登录 非法的用户", seq, request.UserID)
		return
	}
	if !InAppIDs(request.AppID) {
		code = common.Unauthorized
		fmt.Println("用户登录 不支持的平台", seq, request.AppID)
		return
	}
	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.AppID, client.UserID, seq)
		code = common.OperationFailure
		return
	}
	client.Login(request.AppID, request.UserID, currentTime)

	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, request.AppID, request.UserID, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)
		return
	}

	// 用户登录
	login := &login{
		AppID:  request.AppID,
		UserID: request.UserID,
		Client: client,
	}
	clientManager.Login <- login
	fmt.Println("用户登录 成功", seq, client.Addr, request.UserID)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)
		return
	}
	fmt.Println("webSocket_request 心跳接口", client.AppID, client.UserID)
	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.AppID, client.UserID, seq)
		code = common.NotLoggedIn
		return
	}
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			code = common.NotLoggedIn
			fmt.Println("心跳接口 用户未登录", seq, client.AppID, client.UserID)
			return
		} else {
			code = common.ServerError
			fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.AppID, client.UserID, err)
			return
		}
	}
	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppID, client.UserID, err)
		return
	}
	return
}
