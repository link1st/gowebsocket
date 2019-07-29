/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 16:24
 */

package websocket

import (
	"fmt"
	"gowebsocket/lib/cache"
	"sync"
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
func getUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)

	return
}

/**************************  manager  ***************************************/

// 获取用户的连接
func (manager *ClientManager) GetClient(appId uint32, userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := getUserKey(appId, userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}

	return
}

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
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

/**************************  manager event  ***************************************/

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
		userKey := getUserKey(Login.AppId, Login.UserId)
		manager.AddUsers(userKey, Login.Client)
	}

	fmt.Println("EventLogin 用户登录", client.Addr, Login.AppId, Login.UserId)
}

// 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	userKey := getUserKey(client.AppId, client.UserId)
	manager.DelUsers(userKey)

	// 清除redis登录数据
	userOnline, err := cache.GerUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.LogOut()
		cache.SetUserOnlineInfo(client.GetKey(), userOnline)

		// 退出登录事件
		// currentTime := uint64(time.Now().Unix())
	}

	// 关闭 chan
	// close(client.Send)

	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppId, client.UserId)
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
					delete(manager.Clients, conn)
				}
			}
		}

	}
}
