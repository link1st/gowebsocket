// Package websocket 处理
package websocket

import (
	"fmt"
	"sync"
	"time"

	"github.com/link1st/gowebsocket/v2/helper"
	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Login       chan *login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// NewClientManager 创建连接管理
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return
}

// GetUserKey 获取用户key
func GetUserKey(appID uint32, userID string) (key string) {
	key = fmt.Sprintf("%d_%s", appID, userID)
	return
}

func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]
	return
}

// GetClients 获取所有客户端
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

// GetClientsLen GetClientsLen
func (manager *ClientManager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

// AddClients 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// GetUserClient 获取用户的连接
func (manager *ClientManager) GetUserClient(appID uint32, userID string) (client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	userKey := GetUserKey(appID, userID)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}
	return
}

// GetUsersLen GetClientsLen
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

// AddUsers 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	key := GetUserKey(client.AppID, client.UserID)
	if value, ok := manager.Users[key]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, key)
		result = true
	}
	return
}

// GetUserKeys 获取用户的key
func (manager *ClientManager) GetUserKeys() (userKeys []string) {
	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}
	return
}

// GetUserList 获取用户 list
func (manager *ClientManager) GetUserList(appID uint32) (userList []string) {
	userList = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		if v.AppID == appID {
			userList = append(userList, v.UserID)
		}
	}
	fmt.Println("GetUserList len:", len(manager.Users))
	return
}

// GetUserClients 获取用户的key
func (manager *ClientManager) GetUserClients() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// sendAll 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, ignoreClient *Client) {
	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient {
			conn.SendMsg(message)
		}
	}
}

// sendAppIDAll 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAppIDAll(message []byte, appID uint32, ignoreClient *Client) {
	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient && conn.AppID == appID {
			conn.SendMsg(message)
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	fmt.Println("EventRegister 用户建立连接", client.Addr)
	// client.Send <- []byte("连接成功")
}

// EventLogin 用户登录
func (manager *ClientManager) EventLogin(login *login) {
	client := login.Client
	// 连接存在，在添加
	if manager.InClient(client) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}
	fmt.Println("EventLogin 用户登录", client.Addr, login.AppID, login.UserID)
	orderID := helper.GetOrderIDTime()
	_, _ = SendUserMessageAll(login.AppID, login.UserID, orderID, models.MessageCmdEnter, "哈喽~")
}

// EventUnregister 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// 不是当前连接的客户端
		return
	}

	// 清除redis登录数据
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.LogOut()
		_ = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	}

	// 关闭 chan
	// close(client.Send)
	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppID, client.UserID)
	if client.UserID != "" {
		orderID := helper.GetOrderIDTime()
		_, _ = SendUserMessageAll(client.AppID, client.UserID, orderID, models.MessageCmdExit, "用户已经离开~")
	}
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)
		case l := <-manager.Login:
			// 用户登录
			manager.EventLogin(l)
		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)
		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

// GetManagerInfo 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})
	managerInfo["clientsLen"] = clientManager.GetClientsLen()        // 客户端连接数
	managerInfo["usersLen"] = clientManager.GetUsersLen()            // 登录用户数
	managerInfo["chanRegisterLen"] = len(clientManager.Register)     // 未处理连接事件数
	managerInfo["chanLoginLen"] = len(clientManager.Login)           // 未处理登录事件数
	managerInfo["chanUnregisterLen"] = len(clientManager.Unregister) // 未处理退出登录事件数
	managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // 未处理广播事件数
	if isDebug == "true" {
		addrList := make([]string, 0)
		clientManager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)
			return true
		})
		users := clientManager.GetUserKeys()
		managerInfo["clients"] = addrList // 客户端列表
		managerInfo["users"] = users      // 登录用户列表
	}
	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appID uint32, userID string) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())
	clients := clientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserID, client.LoginTime, client.HeartbeatTime)
			_ = client.Socket.Close()
		}
	}
}

// GetUserList 获取全部用户
func GetUserList(appID uint32) (userList []string) {
	fmt.Println("获取全部用户", appID)
	userList = clientManager.GetUserList(appID)
	return
}

// AllSendMessages 全员广播
func AllSendMessages(appID uint32, userID string, data string) {
	fmt.Println("全员广播", appID, userID, data)
	ignoreClient := clientManager.GetUserClient(appID, userID)
	clientManager.sendAppIDAll([]byte(data), appID, ignoreClient)
}
