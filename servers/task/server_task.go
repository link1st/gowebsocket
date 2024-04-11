// Package task 定时任务
package task

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/servers/websocket"
)

// ServerInit 服务初始化
func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// server 服务注册
func server(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务注册 stop", r, string(debug.Stack()))
		}
	}()
	s := websocket.GetServer()
	currentTime := uint64(time.Now().Unix())
	fmt.Println("定时任务，服务注册", param, s, currentTime)
	_ = cache.SetServerInfo(s, currentTime)
	return
}

// serverDefer 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务下线 stop", r, string(debug.Stack()))
		}
	}()
	fmt.Println("服务下线", param)
	s := websocket.GetServer()
	_ = cache.DelServerInfo(s)
	return
}
