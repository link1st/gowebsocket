// Package websocket 处理
package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/models"
)

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
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

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {
	fmt.Println("处理数据", client.Addr, string(message))
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()
	request := &models.Request{}
	if err := json.Unmarshal(message, request); err != nil {
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
	var (
		code uint32
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
	responseHead := models.NewResponseHead(seq, cmd, code, msg, data)
	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		return
	}
	client.SendMsg(headByte)
	fmt.Println("acc_response send", client.Addr, client.AppID, client.UserID, "cmd", cmd, "code", code)
	return
}
