/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-30
* Time: 12:27
 */

package users

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/websocket"
)

// 查询所有用户
func UserList() (userList []string) {

	userList = websocket.GetUserList()

	return
}

// 查询用户是否在线
func CheckUserOnline(appId uint32, userId string) (online bool) {
	// 全平台查询
	if appId == 0 {
		for _, appId := range websocket.GetAppIds() {
			online, _ = checkUserOnline(appId, userId)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appId, userId)
	}

	return
}

// 查询用户 是否在线
func checkUserOnline(appId uint32, userId string) (online bool, err error) {
	key := websocket.GetUserKey(appId, userId)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", appId, userId, err)

			return false, nil
		}

		fmt.Println("GetUserOnlineInfo", appId, userId, err)

		return
	}

	online = userOnline.IsOnline()

	return
}

// 给用户发送消息
func SendUserMessage(appId uint32, userId string, msgId, message string) (sendResults bool, err error) {

	data := models.GetTextMsgData(userId, msgId, message)

	// TODO::需要判断不在本机的情况
	sendResults, err = sendUserMessageLocal(appId, userId, data)
	if err != nil {
		fmt.Println("给用户发送消息", appId, userId, err)
	}

	return
}

// 给本机用户发送消息
func sendUserMessageLocal(appId uint32, userId string, data string) (sendResults bool, err error) {

	client := websocket.GetUserClient(appId, userId)
	if client == nil {
		err = errors.New("用户不在线")

		return
	}

	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true

	return
}

// 给全体用户发消息
func SendUserMessageAll(appId uint32, userId string, msgId, message string) (sendResults bool, err error) {
	sendResults = true

	data := models.GetTextMsgData(userId, msgId, message)
	websocket.AllSendMessages(appId, userId, data)

	return
}
