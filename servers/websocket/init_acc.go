// Package websocket 处理
package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/link1st/gowebsocket/v2/helper"
	"github.com/link1st/gowebsocket/v2/models"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	defaultAppID = 101 // 默认平台ID
)

var (
	clientManager = NewClientManager()                    // 管理者
	appIDs        = []uint32{defaultAppID, 102, 103, 104} // 全部的平台
	serverIp      string
	serverPort    string
)

// GetAppIDs 所有平台
func GetAppIDs() []uint32 {
	return appIDs
}

// GetServer 获取服务器
func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)
	return
}

// IsLocal 判断是否为本机
func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}
	return
}

// InAppIDs in app
func InAppIDs(appID uint32) (inAppID bool) {
	for _, value := range appIDs {
		if value == appID {
			inAppID = true
			return
		}
	}
	return
}

// GetDefaultAppID 获取默认 appID
func GetDefaultAppID() (appID uint32) {
	appID = defaultAppID
	return
}

// StartWebSocket 启动程序
func StartWebSocket() {
	serverIp = helper.GetServerIp()
	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")
	serverPort = rpcPort
	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()
	fmt.Println("WebSocket 启动程序成功", serverIp, serverPort)
	_ = http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)
	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
