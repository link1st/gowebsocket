/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 16:24
 */

package websocket

import (
	"fmt"
	"gowebsocket/helper"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"sync"
	"time"
)

// 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Login       chan *Login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *Login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// 获取用户key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)

	return
}

/**************************  manager  ***************************************/

// 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	delete(manager.Clients, client)
}

// 获取用户的连接
func (manager *ClientManager) GetUserClient(appId uint32, userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := GetUserKey(appId, userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}

	return
}

// 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[key] = client
}

// 删除用户
func (manager *ClientManager) DelUsers(key string) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	delete(manager.Users, key)
}

// 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.SendMsg(message)
		}
	}
}

// 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)

	fmt.Println("EventRegister 用户建立连接", client.Addr)

	// client.Send <- []byte("连接成功")
}

// 用户登录
func (manager *ClientManager) EventLogin(Login *Login) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	client := Login.Client
	// 连接存在，在添加
	if _, ok := manager.Clients[Login.Client]; ok {
		userKey := Login.GetKey()
		manager.AddUsers(userKey, Login.Client)
	}

	fmt.Println("EventLogin 用户登录", client.Addr, Login.AppId, Login.UserId)

	AllSendMessages(Login.AppId, Login.UserId, models.GetTextMsgDataEnter(Login.UserId, helper.GetOrderIdTime(), "哈喽~"))
}

// 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	userKey := GetUserKey(client.AppId, client.UserId)
	manager.DelUsers(userKey)

	// 清除redis登录数据
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.LogOut()
		cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	}

	// 关闭 chan
	// close(client.Send)

	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppId, client.UserId)

	if client.UserId != "" {
		AllSendMessages(client.AppId, client.UserId, models.GetTextMsgDataExit(client.UserId, helper.GetOrderIdTime(), "用户已经离开~"))
	}
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case login := <-manager.Login:
			// 用户登录
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			// 广播事件
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************  manager info  ***************************************/
// 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = len(clientManager.Clients)
	managerInfo["usersLen"] = len(clientManager.Users)
	managerInfo["chanRegisterLen"] = len(clientManager.Register)
	managerInfo["chanLoginLen"] = len(clientManager.Login)
	managerInfo["chanUnregisterLen"] = len(clientManager.Unregister)
	managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)

	if isDebug == "true" {
		clients := make([]string, 0)
		for client := range clientManager.Clients {
			clients = append(clients, client.Addr)
		}

		users := make([]string, 0)
		for key := range clientManager.Users {
			users = append(users, key)
		}

		managerInfo["clients"] = clients
		managerInfo["users"] = users
	}

	return
}

// 获取用户所在的连接
func GetUserClient(appId uint32, userId string) (client *Client) {
	client = clientManager.GetUserClient(appId, userId)

	return
}

// 定时清理超时连接
func ClearTimeoutConnections() {

	currentTime := uint64(time.Now().Unix())

	for client := range clientManager.Clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.LoginTime, client.HeartbeatTime)

			client.Socket.Close()
		}
	}
}

// 获取全部用户
func GetUserList() (userList []string) {

	userList = make([]string, 0)
	fmt.Println("获取全部用户")

	for _, v := range clientManager.Users {
		userList = append(userList, v.UserId)
	}

	return
}

// 全员广播
func AllSendMessages(appId uint32, userId string, data string) {
	fmt.Println("全员广播", appId, userId, data)

	ignore := clientManager.GetUserClient(appId, userId)
	clientManager.sendAll([]byte(data), ignore)
}
