/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-30
* Time: 12:27
 */

package websocket

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/grpcclient"
	"time"
)

// 查询所有用户
func UserList(appID, roomID uint32) (userList []string) {
	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(appID)
		} else {
			list, _ = grpcclient.GetUserList(server, appID, roomID)
		}
		userList = append(userList, list...)
	}
	return
}

// 查询用户是否在线
func CheckUserOnline(appID uint32, userID string) (online bool) {
	// 全平台查询
	if appID == 0 {
		for _, appID := range GetAppIDs() {
			online, _ = checkUserOnline(appID, userID)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appID, userID)
	}
	return
}

// 查询用户 是否在线
func checkUserOnline(appID uint32, userID string) (online bool, err error) {
	key := GetUserKey(appID, userID)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", appID, userID, err)
			return false, nil
		}
		fmt.Println("GetUserOnlineInfo", appID, userID, err)
		return
	}
	online = userOnline.IsOnline()
	return
}

// 给用户发送消息
func SendUserMessage(appID uint32, userID string, msgID, message string) (sendResults bool, err error) {
	data := models.GetTextMsgData(userID, msgID, message)
	client := GetUserClient(appID, userID)
	if client != nil {
		// 在本机发送
		sendResults, err = SendUserMessageLocal(appID, userID, data)
		if err != nil {
			fmt.Println("给用户发送消息", appID, userID, err)
		}
		return
	}
	key := GetUserKey(appID, userID)
	info, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		fmt.Println("给用户发送消息失败", key, err)
		return false, nil
	}
	server := models.NewServer(info.AccIP, info.AccPort)
	msg, err := grpcclient.SendMsg(server, msgID, appID, userID, models.MessageCmdMsg, models.MessageCmdMsg, message)
	if err != nil {
		fmt.Println("给用户发送消息失败", key, err)
		return false, err
	}
	fmt.Println("给用户发送消息成功-rpc", msg)
	sendResults = true
	return
}

// 给本机用户发送消息
func SendUserMessageLocal(appID uint32, userID string, data string) (sendResults bool, err error) {
	client := GetUserClient(appID, userID)
	if client == nil {
		err = errors.New("用户不在线")
		return
	}
	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true
	return
}

// SendRoomMsg 给房间用户发送消息
func SendRoomMsg(appID, roomID uint32, userID string, msgID, cmd, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		if IsLocal(server) {
			data := models.GetMsgData(userID, msgID, cmd, message)
			AllSendMessages(appID, userID, data)
		} else {
			grpcclient.SendRoomMsg(server, msgID, appID, roomID, userID, cmd, message)
		}
	}
	return
}
