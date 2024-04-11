// Package task 定时任务
package task

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/link1st/gowebsocket/v2/servers/websocket"
)

// Init 初始化
func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)

}

// cleanConnection 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()
	fmt.Println("定时任务，清理超时连接", param)
	websocket.ClearTimeoutConnections()
	return
}
