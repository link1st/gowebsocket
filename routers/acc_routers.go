// Package routers 路由
package routers

import (
	"github.com/link1st/gowebsocket/v2/servers/websocket"
)

// WebsocketInit Websocket 路由
func WebsocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("heartbeat", websocket.HeartbeatController)
	websocket.Register("ping", websocket.PingController)
}
