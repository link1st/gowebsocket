# 基于websocket单台机器支持百万连接分布式聊天(IM)系统



本文将介绍如何实现一个基于websocket分布式聊天(IM)系统。

使用golang实现websocket通讯，单机可以支持百万连接，使用gin框架、nginx负载、可以水平部署、程序内部相互通讯、使用grpc通讯协议。

本文内容比较长，如果直接想clone项目体验直接进入[项目体验](#12-项目体验) [goWebSocket项目下载](#4goWebSocket-项目) ,文本从介绍webSocket是什么开始，然后开始介绍这个项目，以及在Nginx中配置域名做webSocket的转发，然后介绍如何搭建一个分布式系统。


## 目录
- [1、项目说明](#1项目说明)
    - [1.1 goWebSocket](#11-goWebSocket)
    - [1.2 项目体验](#12-项目体验)
- [2、介绍webSocket](#2介绍webSocket)
    - [2.1 webSocket 是什么](#21-webSocket-是什么)
    - [2.2 webSocket的兼容性](#22-webSocket的兼容性)
    - [2.3 为什么要用webSocket](#23-为什么要用webSocket)
    - [2.4 webSocket建立过程](#24-webSocket建立过程)
- [3、如何实现基于webSocket的长连接系统](#3如何实现基于webSocket的长连接系统)
    - [3.1 使用go实现webSocket服务端](#31-使用go实现webSocket服务端)
        - [3.1.1 启动端口监听](#311-启动端口监听)
        - [3.1.2 升级协议](#312-升级协议)
        - [3.1.3 客户端连接的管理](#313-客户端连接的管理)
        - [3.1.4 注册客户端的socket的写的异步处理程序](#314-注册客户端的socket的写的异步处理程序)
        - [3.1.5 注册客户端的socket的读的异步处理程序](#315-注册客户端的socket的读的异步处理程序)
        - [3.1.6 接收客户端数据并处理](#316-接收客户端数据并处理)
        - [3.1.7 使用路由的方式处理客户端的请求数据](#317-使用路由的方式处理客户端的请求数据)
        - [3.1.8 防止内存溢出和Goroutine不回收](#318-防止内存溢出和Goroutine不回收)
    - [3.2 使用javaScript实现webSocket客户端](#32-使用javaScript实现webSocket客户端)
        - [3.2.1 启动并注册监听程序](#321-启动并注册监听程序)
        - [3.2.2 发送数据](#322-发送数据)
    - [3.3 发送消息](#33-发送消息)
        - [3.3.1 文本消息](#331-文本消息)
        - [3.3.2 图片和语言消息](#332-图片和语言消息)
- [4、goWebSocket 项目](#4goWebSocket-项目)
    - [4.1 项目说明](#41-项目说明)
    - [4.2 项目依赖](#42-项目依赖)
    - [4.3 项目启动](#43-项目启动)
    - [4.4 接口文档](#44-接口文档)
        - [4.4.1 HTTP接口文档](#441-HTTP接口文档)
            - [4.4.1.1 接口说明](#4411-接口说明)
            - [4.4.1.2 聊天页面](#4412-聊天页面)
            - [4.4.1.3 获取房间用户列表](#4413-获取房间用户列表)
            - [4.4.1.4 查询用户是否在线](#4414-查询用户是否在线)
            - [4.4.1.5 给用户发送消息](#4415-给用户发送消息)
            - [4.4.1.6 给全员用户发送消息](#4416-给全员用户发送消息)
        - [4.4.2 RPC接口文档](#442-RPC接口文档)
            - [4.4.2.1 接口说明](#4421-接口说明)
            - [4.4.2.2 查询用户是否在线](#4422-查询用户是否在线)
            - [4.4.2.3 发送消息](#4423-发送消息)
            - [4.4.2.4 给指定房间所有用户发送消息](#4424-给指定房间所有用户发送消息)
            - [4.4.2.5 获取房间内全部用户](#4425-获取房间内全部用户)
- [5、webSocket项目Nginx配置](#5webSocket项目Nginx配置)
    - [5.1 为什么要配置Nginx](#51-为什么要配置Nginx)
    - [5.2 nginx配置](#52-nginx配置)
    - [5.3 问题处理](#53-问题处理)
- [6、压测](#6压测)
    - [6.1 Linux内核优化](#61-Linux内核优化)
    - [6.2 压测准备](#62-压测准备)
    - [6.3 压测数据](#63-压测数据)
- [7、如何基于webSocket实现一个分布式Im](#7如何基于webSocket实现一个分布式Im)
    - [7.1 说明](#71-说明)
    - [7.2 架构](#72-架构)
    - [7.3 分布式系统部署](#73-分布式系统部署)
- [8、回顾和反思](#8回顾和反思)
    - [8.1 在其它系统应用](#81-在其它系统应用)
    - [8.2 需要完善、优化](#82-需要完善优化)
    - [8.3 总结](#83-总结)
- [9、参考文献](#9参考文献)


## 1、项目说明
#### 1.1 goWebSocket

本文将介绍如何实现一个基于websocket聊天(IM)分布式系统。

使用golang实现websocket通讯，单机支持百万连接，使用gin框架、nginx负载、可以水平部署、程序内部相互通讯、使用grpc通讯协议。

- 一般项目中webSocket使用的架构图
![网站架构图](http://img.91vh.com/img/%E7%BD%91%E7%AB%99%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

#### 1.2 项目体验
- [项目地址 gowebsocket](https://github.com/link1st/gowebsocket)
- [IM-聊天首页](http://im.91vh.com/home/index) 或者在新的窗口打开 http://im.91vh.com/home/index
- 打开连接以后进入聊天界面
- 多人群聊可以同时打开两个窗口

## 2、介绍webSocket
### 2.1 webSocket 是什么
WebSocket 协议在2008年诞生，2011年成为国际标准。所有浏览器都已经支持了。

它的最大特点就是，服务器可以主动向客户端推送信息，客户端也可以主动向服务器发送信息，是真正的双向平等对话，属于服务器推送技术的一种。

- HTTP和WebSocket在通讯过程的比较
![HTTP协议和WebSocket比较](http://img.91vh.com/img/HTTP%E5%8D%8F%E8%AE%AE%E5%92%8CWebSocket%E6%AF%94%E8%BE%83.png)

- HTTP和webSocket都支持配置证书，`ws://` 无证书 `wss://` 配置证书的协议标识
![HTTP协议和WebSocket比较](http://img.91vh.com/img/HTTP%E5%8D%8F%E8%AE%AE%E5%92%8CWebSocket%E6%AF%94%E8%BE%83.jpeg)

### 2.2 webSocket的兼容性
- 浏览器的兼容性，开始支持webSocket的版本

![浏览器开始支持webSocket的版本](http://img.91vh.com/img/%E6%B5%8F%E8%A7%88%E5%99%A8%E5%BC%80%E5%A7%8B%E6%94%AF%E6%8C%81webSocket%E7%9A%84%E7%89%88%E6%9C%AC.jpeg)

- 服务端的支持

golang、java、php、node.js、python、nginx 都有不错的支持

- Android和IOS的支持

Android可以使用java-webSocket对webSocket支持

iOS 4.2及更高版本具有WebSockets支持

### 2.3 为什么要用webSocket
- 1. 从业务上出发，需要一个主动通达客户端的能力
> 目前大多数的请求都是使用HTTP，都是由客户端发起一个请求，有服务端处理，然后返回结果，不可以服务端主动向某一个客户端主动发送数据 

![服务端处理一个请求](http://img.91vh.com/img/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%A4%84%E7%90%86%E4%B8%80%E4%B8%AA%E8%AF%B7%E6%B1%82.jpeg)
- 2. 大多数场景我们需要主动通知用户，如:聊天系统、用户完成任务主动告诉用户、一些运营活动需要通知到在线的用户
- 3. 可以获取用户在线状态
- 4. 在没有长连接的时候通过客户端主动轮询获取数据
- 5. 可以通过一种方式实现，多种不同平台(H5/Android/IOS)去使用

### 2.4 webSocket建立过程
- 1. 客户端先发起升级协议的请求

客户端发起升级协议的请求，采用标准的HTTP报文格式，在报文中添加头部信息 

`Connection: Upgrade`表明连接需要升级

`Upgrade: websocket`需要升级到 websocket协议

`Sec-WebSocket-Version: 13` 协议的版本为13

`Sec-WebSocket-Key: I6qjdEaqYljv3+9x+GrhqA==` 这个是base64 encode 的值，是浏览器随机生成的，与服务器响应的 `Sec-WebSocket-Accept`对应

```
# Request Headers
Connection: Upgrade
Host: im.91vh.com
Origin: http://im.91vh.com
Pragma: no-cache
Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
Sec-WebSocket-Key: I6qjdEaqYljv3+9x+GrhqA==
Sec-WebSocket-Version: 13
Upgrade: websocket
```

![浏览器 Network](http://img.91vh.com/img/%E6%B5%8F%E8%A7%88%E5%99%A8%20Network.png)

- 2. 服务器响应升级协议

服务端接收到升级协议的请求，如果服务端支持升级协议会做如下响应

返回: 

`Status Code: 101 Switching Protocols` 表示支持切换协议

```
# Response Headers
Connection: upgrade
Date: Fri, 09 Aug 2019 07:36:59 GMT
Sec-WebSocket-Accept: mB5emvxi2jwTUhDdlRtADuBax9E=
Server: nginx/1.12.1
Upgrade: websocket
```

- 3. 升级协议完成以后，客户端和服务器就可以相互发送数据

![websocket接收和发送数据](http://img.91vh.com/img/websocket%E6%8E%A5%E6%94%B6%E5%92%8C%E5%8F%91%E9%80%81%E6%95%B0%E6%8D%AE.png)

## 3、如何实现基于webSocket的长连接系统

### 3.1 使用go实现webSocket服务端

#### 3.1.1 启动端口监听
- websocket需要监听端口，所以需要在`golang` 程序的 `main` 函数中用协程的方式去启动程序
- **main.go** 实现启动

```
go websocket.StartWebSocket()
```
- **init_acc.go** 启动程序

```
// 启动程序
func StartWebSocket() {
	http.HandleFunc("/acc", wsPage)
	http.ListenAndServe(":8089", nil)
}
```

#### 3.1.2 升级协议
- 客户端是通过http请求发送到服务端，我们需要对http协议进行升级为websocket协议
- 对http请求协议进行升级 golang 库[gorilla/websocket](https://github.com/gorilla/websocket) 已经做得很好了，我们直接使用就可以了
- 在实际使用的时候，建议每个连接使用两个协程处理客户端请求数据和向客户端发送数据，虽然开启协程会占用一些内存，但是读取分离，减少收发数据堵塞的可能
- **init_acc.go**

```
func wsPage(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])

		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)

		return
	}

	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
```

#### 3.1.3 客户端连接的管理
- 当前程序有多少用户连接，还需要对用户广播的需要，这里我们就需要一个管理者(clientManager)，处理这些事件:
- 记录全部的连接、登录用户的可以通过 **appId+uuid** 查到用户连接
- 使用map存储，就涉及到多协程并发读写的问题，所以需要加读写锁
- 定义四个channel ，分别处理客户端建立连接、用户登录、断开连接、全员广播事件

```
// 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Login       chan *login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// 初始化
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
```

#### 3.1.4 注册客户端的socket的写的异步处理程序
- 防止发生程序崩溃，所以需要捕获异常
- 为了显示异常崩溃位置这里使用`string(debug.Stack())`打印调用堆栈信息
- 如果写入数据失败了，可能连接有问题，就关闭连接
- **client.go**

```
// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)

		}
	}()

	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
		fmt.Println("Client发送数据 defer", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				fmt.Println("Client发送数据 关闭连接", c.Addr, "ok", ok)

				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
```

#### 3.1.5 注册客户端的socket的读的异步处理程序
- 循环读取客户端发送的数据并处理
- 如果读取数据失败了，关闭channel
- **client.go**

```
// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		fmt.Println("读取客户端数据 关闭send", c)
		close(c.Send)
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端数据 错误", c.Addr, err)

			return
		}

		// 处理程序
		fmt.Println("读取客户端数据 处理:", string(message))
		ProcessData(c, message)
	}
}
```

#### 3.1.6 接收客户端数据并处理
- 约定发送和接收请求数据格式，为了js处理方便，采用了`json`的数据格式发送和接收数据(人类可以阅读的格式在工作开发中使用是比较方便的)

- 登录发送数据示例:
```
{"seq":"1565336219141-266129","cmd":"login","data":{"userId":"马远","appId":101}}
```
- 登录响应数据示例:
```
{"seq":"1565336219141-266129","cmd":"login","response":{"code":200,"codeMsg":"Success","data":null}}
```
- websocket是双向的数据通讯，可以连续发送，如果发送的数据需要服务端回复，就需要一个**seq**来确定服务端的响应是回复哪一次的请求数据
- cmd 是用来确定动作，websocket没有类似于http的url,所以规定 cmd 是什么动作
- 目前的动作有:login/heartbeat 用来发送登录请求和连接保活(长时间没有数据发送的长连接容易被浏览器、移动中间商、nginx、服务端程序断开)
- 为什么需要AppId,UserId是表示用户的唯一字段，设计的时候为了做成通用性，设计AppId用来表示用户在哪个平台登录的(web、app、ios等)，方便后续扩展

- **request_model.go** 约定的请求数据格式

```
/************************  请求数据  **************************/
// 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
}

// 心跳请求数据
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}
```

- **response_model.go**

```
/************************  响应数据  **************************/
type Head struct {
	Seq      string    `json:"seq"`      // 消息的Id
	Cmd      string    `json:"cmd"`      // 消息的cmd 动作
	Response *Response `json:"response"` // 消息体
}

type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

```


#### 3.1.7 使用路由的方式处理客户端的请求数据

- 使用路由的方式处理由客户端发送过来的请求数据
- 以后添加请求类型以后就可以用类是用http相类似的方式(router-controller)去处理
- **acc_routers.go**

```
// Websocket 路由
func WebsocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("heartbeat", websocket.HeartbeatController)
}
```

#### 3.1.8 防止内存溢出和Goroutine不回收
- 1. 定时任务清除超时连接
没有登录的连接和登录的连接6分钟没有心跳则断开连接

**client_manager.go**

```
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
```

- 2. 读写的Goroutine有一个失败，则相互关闭
`write()`Goroutine写入数据失败，关闭`c.Socket.Close()`连接，会关闭`read()`Goroutine
`read()`Goroutine读取数据失败，关闭`close(c.Send)`连接，会关闭`write()`Goroutine

- 3. 客户端主动关闭
关闭读写的Goroutine
从`ClientManager`删除连接

- 4. 监控用户连接、Goroutine数
十个内存溢出有九个和Goroutine有关
添加一个http的接口，可以查看系统的状态，防止Goroutine不回收
[查看系统状态](http://im.91vh.com/system/state?isDebug=true)

- 5. Nginx 配置不活跃的连接释放时间，防止忘记关闭的连接

- 6. 使用 pprof 分析性能、耗时

### 3.2 使用javaScript实现webSocket客户端
#### 3.2.1 启动并注册监听程序
- js 建立连接，并处理连接成功、收到数据、断开连接的事件处理

```
ws = new WebSocket("ws://127.0.0.1:8089/acc");

 
ws.onopen = function(evt) {
  console.log("Connection open ...");
};
 
ws.onmessage = function(evt) {
  console.log( "Received Message: " + evt.data);
  data_array = JSON.parse(evt.data);
  console.log( data_array);
};
 
ws.onclose = function(evt) {
  console.log("Connection closed.");
};

```


#### 3.2.2 发送数据
- 需要注意:连接建立成功以后才可以发送数据
- 建立连接以后由客户端向服务器发送数据示例

```
登录:
ws.send('{"seq":"2323","cmd":"login","data":{"userId":"11","appId":101}}');

心跳:
ws.send('{"seq":"2324","cmd":"heartbeat","data":{}}');

ping 查看服务是否正常:
ws.send('{"seq":"2325","cmd":"ping","data":{}}');

关闭连接:
ws.close();
```

## 3.3 发送消息
### 3.3.1 文本消息

客户端只要知道发送用户是谁，还有内容就可以显示文本消息，这里我们重点关注一下数据部分

target：定义接收的目标，目前未设置

type：消息的类型，text 文本消息 img 图片消息 

msg：文本消息内容

from：消息的发送者

文本消息的结构:

```json
{
  "seq": "1569080188418-747717",
  "cmd": "msg",
  "response": {
    "code": 200,
    "codeMsg": "Ok",
    "data": {
      "target": "",
      "type": "text",
      "msg": "hello",
      "from": "马超"
    }
  }
}
```

这样一个文本消息的结构就设计完成了，客户端在接收到消息内容就可以展现到 IM 界面上

### 3.3.2 图片和语言消息

发送图片消息，发送消息者的客户端需要先把图片上传到文件服务器，上传成功以后获得图片访问的 URL，然后由发送消息者的客户端需要将图片 URL 发送到 gowebsocket，gowebsocket 图片的消息格式发送给目标客户端，消息接收者客户端接收到图片的 URL 就可以显示图片消息。

图片消息的结构:

```
{
  "type": "img",
  "from": "马超",
  "url": "http://91vh.com/images/home_logo.png",
  "secret": "消息鉴权 secret",
  "size": {
    "width": 480,
    "height": 720
  }
}
```

语言消息、和视频消息和图片消息类似，都是先把文件上传服务器，然后通过 gowebsocket 传递文件的 URL，需要注意的是部分消息涉及到隐私的文件，文件访问的时候需要做好鉴权信息，不能让非接收用户也能查看到别人的消息内容。

## 4、goWebSocket 项目
### 4.1 项目说明
- 本项目是基于webSocket实现的分布式IM系统
- 客户端随机分配用户名，所有人进入一个聊天室，实现群聊的功能
- 单台机器(24核128G内存)支持百万客户端连接
- 支持水平部署，部署的机器之间可以相互通讯

- 项目架构图
![网站架构图](http://img.91vh.com/img/%E7%BD%91%E7%AB%99%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

### 4.2 项目依赖

- 本项目只需要使用 redis 和 golang 
- 本项目使用govendor管理依赖，克隆本项目就可以直接使用

```
# 主要使用到的包
github.com/gin-gonic/gin@v1.4.0
github.com/redis/go-redis/v9
github.com/gorilla/websocket
github.com/spf13/viper
google.golang.org/grpc
github.com/golang/protobuf
```


### 4.3 项目启动 
- 克隆项目

```
git clone git@github.com:link1st/gowebsocket.git
# 或
git clone https://github.com/link1st/gowebsocket.git
```
- 修改项目配置

```
cd gowebsocket
cd config
mv app.yaml.example app.yaml
# 修改项目监听端口，redis连接等(默认127.0.0.1:3306)
vim app.yaml
# 返回项目目录，为以后启动做准备
cd ..
```
- 配置文件说明

```
app:
  logFile: log/gin.log # 日志文件位置
  httpPort: 8080 # http端口
  webSocketPort: 8089 # webSocket端口
  rpcPort: 9001 # 分布式部署程序内部通讯端口
  httpUrl: 127.0.0.1:8080
  webSocketUrl:  127.0.0.1:8089


redis:
  addr: "localhost:6379"
  password: ""
  DB: 0
  poolSize: 30
  minIdleConns: 30
```

- 启动项目

```
go run main.go
```

- 进入IM聊天地址
[http://127.0.0.1:8080/home/index](http://127.0.0.1:8080/home/index)
- 到这里，就可以体验到基于webSocket的IM系统

#### 4.4 接口文档 
###### 4.4.1.1 接口说明
##### 4.4.1 HTTP接口文档
- 在接口开发和接口文档使用的过程中，规范开发流程，减少沟通成本，所以约定一下接口开发流程和文档说明
- 接口地址

 线上:http://im.91vh.com

 测试:http://im.91vh.com


###### 4.4.1.2 聊天页面
- 地址:/home/index
- 请求方式:GET
- 接口说明:聊天页面
- 请求参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| appId   |   是    | uint32 | appId/房间Id |   101      |

- 返回参数:
无


###### 4.4.1.3 获取房间用户列表
- 地址:/user/list
- 请求方式:GET/POST
- 接口说明:获取房间用户列表
- 请求参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| appId   |   是    | uint32 | appId/房间Id |   101      |

- 返回参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| code   |   是    | int   | 错误码  |   200  |
| msg    |   是    | string| 错误信息 |Success |
| data   |   是    | array | 返回数据 |        |
| userCount   |   是    | int   | 房间内用户总数  |   1    |
| userList| 是 | list  | 用户列表 |        |

- 示例:

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "userCount": 1,
        "userList": [
            "黄帝"
        ]
    }
}
```

###### 4.4.1.4 查询用户是否在线
- 地址:/user/online
- 请求方式:GET/POST
- 接口说明:查询用户是否在线
- 请求参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| appId   |   是    | uint32 | appId/房间Id |   101      |
| userId   |   是    | string | 用户Id |   黄帝     |

- 返回参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| code   |   是    | int   | 错误码  |   200  |
| msg    |   是    | string| 错误信息 |Success |
| data   |   是    | array | 返回数据 |        |
| online   |   是    | bool   | 发送结果 true:在线 false:不在线  |   true    |
| userId   |   是    | string | 用户Id |   黄帝     |

- 示例:

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "online": true,
        "userId": "黄帝"
    }
}
```

###### 4.4.1.5 给用户发送消息
- 地址:/user/sendMessage
- 请求方式:GET/POST
- 接口说明:给用户发送消息
- 请求参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| appId   |   是    | uint32 | appId/房间Id |   101      |
| userId   |   是    | string | 用户id |   黄帝      |
| msgId   |   是    | string | 消息Id |   避免重复发送      |
| message   |   是    | string | 消息内容 |   hello      |

- 返回参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| code   |   是    | int   | 错误码  |   200  |
| msg    |   是    | string| 错误信息 |Success |
| data   |   是    | array | 返回数据 |        |
| sendResults   |   是    | bool   | 发送结果 true:成功 false:失败  |   true    |

- 示例:

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "sendResults": true
    }
}
```

###### 4.4.1.6 给全员用户发送消息
- 地址:/user/sendMessageAll
- 请求方式:GET/POST
- 接口说明:给全员用户发送消息
- 请求参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| appId   |   是    | uint32 | appId/房间Id |   101      |
| userId   |   是    | string | 用户id |   黄帝      |
| msgId   |   是    | string | 消息Id |   避免重复发送      |
| message   |   是    | string | 消息内容 |   hello      |

- 返回参数:

|  参数   |  必填   |  类型  |  说明   |  示例   |
| :----: | :----: | :----: | :----: | :----: |
| code   |   是    | int   | 错误码  |   200  |
| msg    |   是    | string| 错误信息 |Success |
| data   |   是    | array | 返回数据 |        |
| sendResults   |   是    | bool   | 发送结果 true:成功 false:失败  |   true    |

- 示例:

```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "sendResults": true
    }
}
```

##### 4.4.2 RPC接口文档
###### 4.4.2.1 接口说明
- 接口协议结构体
```proto
syntax = "proto3";

option java_multiple_files = true;
option java_package = "io.grpc.examples.protobuf";
option java_outer_classname = "ProtobufProto";


package protobuf;

// The AccServer service definition.
service AccServer {
    // 查询用户是否在线
    rpc QueryUsersOnline (QueryUsersOnlineReq) returns (QueryUsersOnlineRsp) {
    }
    // 发送消息
    rpc SendMsg (SendMsgReq) returns (SendMsgRsp) {
    }
    // 给这台机器的房间内所有用户发送消息
    rpc SendMsgAll (SendMsgAllReq) returns (SendMsgAllRsp) {
    }
    // 获取用户列表
    rpc GetUserList (GetUserListReq) returns (GetUserListRsp) {
    }
}

// 查询用户是否在线
message QueryUsersOnlineReq {
    uint32 appId = 1; // AppID
    string userId = 2; // 用户ID
}

message QueryUsersOnlineRsp {
    uint32 retCode = 1;
    string errMsg = 2;
    bool online = 3;
}

// 发送消息
message SendMsgReq {
    string seq = 1; // 序列号
    uint32 appId = 2; // appId/房间Id
    string userId = 3; // 用户ID
    string cms = 4; // cms 动作: msg/enter/exit
    string type = 5; // type 消息类型，默认是 text
    string msg = 6; // msg
    bool isLocal = 7; // 是否查询本机 acc内部调用为:true(本机查询不到即结束)
}

message SendMsgRsp {
    uint32 retCode = 1;
    string errMsg = 2;
    string sendMsgId = 3;
}

// 给这台机器的房间内所有用户发送消息
message SendMsgAllReq {
    string seq = 1; // 序列号
    uint32 appId = 2; // appId/房间Id
    string userId = 3; // 不发送的用户ID
    string cms = 4; // cms 动作: msg/enter/exit
    string type = 5; // type 消息类型，默认是 text
    string msg = 6; // msg
}

message SendMsgAllRsp {
    uint32 retCode = 1;
    string errMsg = 2;
    string sendMsgId = 3;
}

// 获取用户列表
message GetUserListReq {
    uint32 appId = 1;
}

message GetUserListRsp {
    uint32 retCode = 1;
    string errMsg = 2;
    repeated string userId = 3;
}
```

###### 4.4.2.2 查询用户是否在线
- 参考上述协议结构体

###### 4.4.2.3 发送消息
###### 4.4.2.4 给指定房间所有用户发送消息
###### 4.4.2.5 获取房间内全部用户

## 5、webSocket项目Nginx配置
### 5.1 为什么要配置Nginx
- 使用nginx实现内外网分离，对外只暴露Nginx的Ip(一般的互联网企业会在nginx之前加一层LVS做负载均衡)，减少入侵的可能
- 支持配置 ssl 证书，使用 `wss` 的方式实现数据加密，减少数据被抓包和篡改的可能性
- 使用Nginx可以利用Nginx的负载功能，前端再使用的时候只需要连接固定的域名，通过Nginx将流量分发了到不同的机器
- 同时我们也可以使用Nginx的不同的负载策略(轮询、weight、ip_hash)

### 5.2 nginx配置
- 使用域名 **im.91vh.com** 为示例，参考配置
- 一级目录**im.91vh.com/acc** 是给webSocket使用，是用nginx stream转发功能(nginx 1.3.31 开始支持，使用Tengine配置也是相同的)，转发到golang 8089 端口处理
- 其它目录是给HTTP使用，转发到golang 8080 端口处理

```
upstream  go-im
{
    server 127.0.0.1:8080 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}

upstream  go-acc
{
    server 127.0.0.1:8089 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}


server {
    listen       80 ;
    server_name  im.91vh.com;
    index index.html index.htm ;


    location /acc {
        proxy_set_header Host $host;
        proxy_pass http://go-acc;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Connection "";
        proxy_redirect off;
        proxy_intercept_errors on;
        client_max_body_size 10m;
    }

    location /
    {
        proxy_set_header Host $host;
        proxy_pass http://go-im;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_redirect off;
        proxy_intercept_errors on;
        client_max_body_size 30m;
    }

    access_log  /link/log/nginx/access/im.log;
    error_log   /link/log/nginx/access/im.error.log;
}
```

### 5.3 问题处理
- 运行nginx测试命令，查看配置文件是否正确

```
/link/server/tengine/sbin/nginx -t

```

- 如果出现错误

```
nginx: [emerg] unknown "connection_upgrade" variable
configuration file /link/server/tengine/conf/nginx.conf test failed
```

- 处理方法
- 在**nginx.com**添加

```
http{
	fastcgi_temp_file_write_size 128k;
..... # 需要添加的内容

    #support websocket
    map $http_upgrade $connection_upgrade {
        default upgrade;
        ''      close;
    }

.....
    gzip on;
    
}

```

- 原因:Nginx代理webSocket的时候就会遇到Nginx的设计问题 **End-to-end and Hop-by-hop Headers** 


## 6、压测
### 6.1 Linux内核优化
- 设置文件打开句柄数

被压测服务器需要保持100W长连接，客户和服务器端是通过socket通讯的，每个连接需要建立一个socket，程序需要保持100W长连接就需要单个程序能打开100W个文件句柄

```
# 查看系统默认的值
ulimit -n
# 设置最大打开文件数
ulimit -n 1000000
```

通过修改配置文件的方式修改程序最大打开句柄数

```
root soft nofile 1040000
root hard nofile 1040000

root soft nofile 1040000
root hard nproc 1040000

root soft core unlimited
root hard core unlimited

* soft nofile 1040000
* hard nofile 1040000

* soft nofile 1040000
* hard nproc 1040000

* soft core unlimited
* hard core unlimited
```

修改完成以后需要重启机器配置才能生效

- 修改系统级别文件句柄数量

file-max的值需要大于limits设置的值

```
# file-max 设置的值参考
cat /proc/sys/fs/file-max
12553500
```

- 设置sockets连接参数

`vim /etc/sysctl.conf` 

```
# 配置参考
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 0
net.ipv4.ip_local_port_range = 1024 65000
net.ipv4.tcp_mem = 786432 2097152 3145728
net.ipv4.tcp_rmem = 4096 4096 16777216
net.ipv4.tcp_wmem = 4096 4096 16777216
```

`sysctl -p` 修改配置以后使得配置生效命令

### 6.2 压测准备
- 待压测，如果大家有压测的结果欢迎补充
- 后续会出专门的教程,从申请机器、写压测用例、内核优化、得出压测数据

- **关于压测请移步**
- [go实现的压测工具【单台机器100w连接压测实战】](https://github.com/link1st/go-stress-testing)
- 用go语言实现一款压测工具，然后对本项目进行压测，实现单台机器100W长连接

### 6.3 压测数据
- 项目在实际使用的时候，每个连接约占 27Kb内存
- 支持百万连接需要25G内存，单台机器实现百万长连接是可以实现的

- 记录内存使用情况，分别记录了1W到100W连接数内存使用情况

| 连接数      |  内存 |
| :----:     | :----:|
|   10000    | 281M  |
|   100000   | 2.7g  |
|   200000   | 5.4g  |
|   500000   | 13.1g |
|   1000000  | 25.8g |

- [压测详细数据](https://github.com/link1st/go-stress-testing#65-%E5%8E%8B%E6%B5%8B%E6%95%B0%E6%8D%AE)

## 7、如何基于webSocket实现一个分布式Im
### 7.1 说明
- 参考本项目源码
- [gowebsocket v1.0.0 单机版Im系统](https://github.com/link1st/gowebsocket/tree/v1.0.0)
- [gowebsocket v2.0.0 分布式Im系统](https://github.com/link1st/gowebsocket/tree/v2.0.0)

- 为了方便演示，IM系统和webSocket(acc)系统合并在一个系统中
- IM系统接口:
获取全部在线的用户，查询当前服务的全部用户+集群中服务的全部用户
发送消息，这里采用的是http接口发送(微信网页版发送消息也是http接口)，这里考虑主要是两点:
1.服务分离，让acc系统尽量的简单一点，不掺杂其它业务逻辑
2.发送消息是走http接口，不使用webSocket连接，采用收和发送数据分离的方式，可以加快收发数据的效率

### 7.2 架构

- 项目启动注册和用户连接时序图

![用户连接时序图](http://img.91vh.com/img/%E7%94%A8%E6%88%B7%E8%BF%9E%E6%8E%A5%E6%97%B6%E5%BA%8F%E5%9B%BE.png)

- 其它系统(IM、任务)向webSocket(acc)系统连接的用户发送消息时序图

![分布是系统随机给用户发送消息](http://img.91vh.com/img/%E5%88%86%E5%B8%83%E6%98%AF%E7%B3%BB%E7%BB%9F%E9%9A%8F%E6%9C%BA%E7%BB%99%E7%94%A8%E6%88%B7%E5%8F%91%E9%80%81%E6%B6%88%E6%81%AF.png)

### 7.3 分布式系统部署
- 用水平部署两个项目(gowebsocket和gowebsocket1)演示分部署
- 项目之间如何相互通讯:项目启动以后将项目Ip、rpcPort注册到redis中，让其它项目可以发现，需要通讯的时候使用gRpc进行通讯
- gowebsocket

```
# app.yaml 配置文件信息
app:
  logFile: log/gin.log
  httpPort: 8080
  webSocketPort: 8089
  rpcPort: 9001
  httpUrl: im.91vh.com
  webSocketUrl:  im.91vh.com

# 在启动项目
go run main.go 

```

- gowebsocket1 

```
# 将第一个项目拷贝一份
cp -rf gowebsocket gowebsocket1
# app.yaml 修改配置文件
app:
  logFile: log/gin.log
  httpPort: 8081
  webSocketPort: 8090
  rpcPort: 9002
  httpUrl: im.91vh.com
  webSocketUrl:  im.91vh.com

# 在启动第二个项目
go run main.go 
```

- Nginx配置

在之前Nginx配置项中添加第二台机器的Ip和端口

```
upstream  go-im
{
    server 127.0.0.1:8080 weight=1 max_fails=2 fail_timeout=10s;
    server 127.0.0.1:8081 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}

upstream  go-acc
{
    server 127.0.0.1:8089 weight=1 max_fails=2 fail_timeout=10s;
    server 127.0.0.1:8090 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}
```

- 配置完成以后重启Nginx
- 重启以后请求，验证是否符合预期:

 查看请求是否落在两个项目上
 实验两个用户分别连接不同的项目(gowebsocket和gowebsocket1)是否也可以相互发送消息

- 关于分布式部署

 本项目只是演示了这个项目如何分布式部署，以及分布式部署以后模块如何进行相互通讯
 完全解决系统没有单点的故障，还需 Nginx集群、redis cluster等


## 8、回顾和反思
### 8.1 在其它系统应用
- 本系统设计的初衷就是:和客户端保持一个长连接、对外部系统两个接口(查询用户是否在线、给在线的用户推送消息)，实现业务的分离
- 只有和业务分离可，才可以供多个业务使用，而不是每个业务都建立一个长连接

#### 8.2 已经实现的功能

- gin log日志(请求日志+debug日志)
- 读取配置文件 完成
- 定时脚本，清理过期未心跳连接 完成
- http接口，获取登录、连接数量 完成
- http接口，发送push、查询有多少人在线 完成
- grpc 程序内部通讯，发送消息 完成
- appIds 一个用户在多个平台登录
- 界面，把所有在线的人拉倒一个群里面，发送消息 完成
- ~~单聊~~、群聊 完成
- 实现分布式，水平扩张 完成
- 压测脚本
- 文档整理
- 文档目录、百万长连接的实现、为什么要实现一个IM、怎么实现一个Im 
- 架构图以及扩展

IM实现细节:

- 定义文本消息结构 完成
- html发送文本消息 完成
- 接口接收文本消息并发送给全体 完成
- html接收到消息 显示到界面 完成
- 界面优化 需要持续优化
- 有人加入以后广播全体 完成
- 定义加入聊天室的消息结构 完成
- 引入机器人 待定

### 8.2 需要完善、优化
- 登录，使用微信登录 获取昵称、头像等
- 有账号系统、资料系统
- 界面优化、适配手机端
- 消息 文本消息(支持表情)、图片、语音、视频消息
- 微服务注册、发现、熔断等
- 添加配置项，单台机器最大连接数量

### 8.3 总结
- 虽然实现了一个分布式在聊天的IM，但是有很多细节没有处理(登录没有鉴权、界面还待优化等)，但是可以通过这个示例可以了解到:通过WebSocket解决很多业务上需求
- 本文虽然号称单台机器能有百万长连接(内存上能满足)，但是实际在场景远比这个复杂(cpu有些压力)，当然了如果你有这么大的业务量可以购买更多的机器更好的去支撑你的业务，本程序只是演示如何在实际工作用使用webSocket.
- 参考本文，你可以实现出来符合你需要的程序

### 9、参考文献

[维基百科 WebSocket](https://zh.wikipedia.org/wiki/WebSocket)

[阮一峰 WebSocket教程](http://www.ruanyifeng.com/blog/2017/05/websocket.html)

[WebSocket协议：5分钟从入门到精通](https://www.cnblogs.com/chyingp/p/websocket-deep-in.html)

[go-stress-testing 单台机器100w连接压测实战](https://github.com/link1st/go-stress-testing)

github 搜:link1st 查看项目 gowebsocket

[https://github.com/link1st/gowebsocket](https://github.com/link1st/gowebsocket)

### 意见反馈

- 在项目中遇到问题可以直接在这里找找答案或者提问 [issues](https://github.com/link1st/gowebsocket/issues)
- 也可以添加我的微信(申请信息填写:公司、姓名，我好备注下)，直接反馈给我
<br/>
<p align="center">
     <img border="0" src="http://img.91vh.com/img/%E5%BE%AE%E4%BF%A1%E4%BA%8C%E7%BB%B4%E7%A0%81.jpeg" alt="添加link1st的微信" width="200"/>
</p>

### 赞助商

- 感谢[JetBrains](https://www.jetbrains.com/?from=gowebsocket)对本项目的支持！
<br/>
<p align="center">
    <a href="https://www.jetbrains.com/?from=gowebsocket">
        <img border="0" src="http://img.91vh.com/img/jetbrains_logo.png" width="200"/>
    </a>
</p>
