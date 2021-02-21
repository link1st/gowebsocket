/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gowebsocket/common"
	"gowebsocket/controllers"
	"gowebsocket/helper"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/websocket"
	"strconv"
)

// List 查看全部在线用户
func List(c *gin.Context) {
	appID := helper.StrToUint32(c.Query("appID"))
	roomID := helper.StrToUint32(c.Query("roomID"))
	fmt.Println("http_request 查看全部在线用户", appID, roomID)
	data := make(map[string]interface{})

	userList := websocket.UserList(appID, roomID)
	data["userList"] = userList
	data["userCount"] = len(userList)

	controllers.Response(c, common.OK, "", data)
}

// 查看用户是否在线
func Online(c *gin.Context) {
	userID := c.Query("userID")
	appID := helper.StrToUint32(c.Query("appID"))
	fmt.Println("http_request 查看用户是否在线", userID, appID)

	data := make(map[string]interface{})
	online := websocket.CheckUserOnline(appID, userID)
	data["userID"] = userID
	data["online"] = online
	controllers.Response(c, common.OK, "", data)
}

// 给用户发送消息
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

// SendRoomMessage 群成员发送消息
func SendRoomMessage(c *gin.Context) {
	// 获取参数
	appID := helper.StrToUint32(c.PostForm("appID"))
	roomID := helper.StrToUint32(c.PostForm("roomID"))
	userID := c.PostForm("userID")
	msgID := c.PostForm("msgID")
	message := c.PostForm("message")

	fmt.Println("http_request 给全体用户发送消息", appID, roomID, userID, msgID, message)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendRoomMsg(appID, roomID, userID, msgID, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)

}
