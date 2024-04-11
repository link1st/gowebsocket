// Package user 用户调用接口
package user

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/controllers"
	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/servers/websocket"
)

// List 查看全部在线用户
func List(c *gin.Context) {
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看全部在线用户", appID)
	data := make(map[string]interface{})
	userList := websocket.UserList(appID)
	data["userList"] = userList
	data["userCount"] = len(userList)
	controllers.Response(c, common.OK, "", data)
}

// Online 查看用户是否在线
func Online(c *gin.Context) {
	userID := c.Query("userID")
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看用户是否在线", userID, appIDStr)
	data := make(map[string]interface{})
	online := websocket.CheckUserOnline(appID, userID)
	data["userID"] = userID
	data["online"] = online
	controllers.Response(c, common.OK, "", data)
}

// SendMessage 给用户发送消息
func SendMessage(c *gin.Context) {
	// 获取参数
	appIDStr := c.PostForm("appID")
	userID := c.PostForm("userID")
	msgID := c.PostForm("msgID")
	message := c.PostForm("message")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 给用户发送消息", appIDStr, userID, msgID, message)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID(userID)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendUserMessage(appID, userID, msgID, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)
}

// SendMessageAll 给全员发送消息
func SendMessageAll(c *gin.Context) {
	// 获取参数
	appIDStr := c.PostForm("appID")
	userID := c.PostForm("userID")
	msgID := c.PostForm("msgID")
	message := c.PostForm("message")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 给全体用户发送消息", appIDStr, userID, msgID, message)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendUserMessageAll(appID, userID, msgID, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)

}
