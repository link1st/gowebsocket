// Package websocket 处理
package websocket

import (
	"errors"
	"fmt"
	"time"

	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/servers/grpcclient"

	"github.com/redis/go-redis/v9"
)

// UserList 查询所有用户
func UserList(appID uint32) (userList []string) {
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
			list, _ = grpcclient.GetUserList(server, appID)
		}
		userList = append(userList, list...)
	}
	return
}

// CheckUserOnline 查询用户是否在线
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

// checkUserOnline 查询用户 是否在线
func checkUserOnline(appID uint32, userID string) (online bool, err error) {
	key := GetUserKey(appID, userID)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("GetUserOnlineInfo", appID, userID, err)
			return false, nil
		}
		fmt.Println("GetUserOnlineInfo", appID, userID, err)
		return
	}
	online = userOnline.IsOnline()
	return
}

// SendUserMessage 给用户发送消息
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
	if !info.IsOnline() {
		fmt.Println("用户不在线", key)
		return false, nil
	}
	server := models.NewServer(info.AccIp, info.AccPort)
	msg, err := grpcclient.SendMsg(server, msgID, appID, userID, models.MessageCmdMsg, models.MessageCmdMsg, message)
	if err != nil {
		fmt.Println("给用户发送消息失败", key, err)
		return false, err
	}
	fmt.Println("给用户发送消息成功-rpc", msg)
	sendResults = true
	return
}

// SendUserMessageLocal 给本机用户发送消息
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

// SendUserMessageAll 给全体用户发消息
func SendUserMessageAll(appID uint32, userID string, msgID, cmd, message string) (sendResults bool, err error) {
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
			_, _ = grpcclient.SendMsgAll(server, msgID, appID, userID, cmd, message)
		}
	}
	return
}
