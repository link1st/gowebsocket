/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package systems

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/link1st/gowebsocket/common"
	"github.com/link1st/gowebsocket/controllers"
	"github.com/link1st/gowebsocket/servers/websocket"
)

// 查询系统状态
func Status(c *gin.Context) {

	isDebug := c.Query("isDebug")
	fmt.Println("http_request 查询系统状态", isDebug)

	data := make(map[string]interface{})

	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager 信息
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)

	controllers.Response(c, common.OK, "", data)
}
