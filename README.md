## 目录
* [1、项目说明](#1项目说明)
    * [1.1 goWebSocket](#11-goWebSocket)
    * [1.2 项目体验](#12-项目体验)
- [2、介绍webSocket](#2介绍webSocket)
    - [2.1 webSocket 是什么](#21-webSocket-是什么)
    - [2.2 webSocket的兼容性](#22-webSocket的兼容性)
    - [2.3 为什么要用webSocket](#23-为什么要用webSocket)
- [3、如何实现基于webSocket的长链接系统](#3如何实现基于webSocket的长链接系统)
    - [3.1 使用go实现webSocket服务端](#31-使用go实现webSocket服务端)
        - [3.1.1 启动端口兼容听](#311-启动端口兼容听)
        - [3.1.2 升级协议](#312-升级协议)
        - [3.1.3 客户端连接的管理](#313-客户端连接的管理)
        - [3.1.4 注册客户端的socket的写的异步处理程序](#314-注册客户端的socket的写的异步处理程序)
        - [3.1.5 注册客户端的socket的读的异步处理程序](#315-注册客户端的socket的读的异步处理程序)
        - [3.1.6 接收客户端数据并处理](#316-接收客户端数据并处理)
        - [3.1.7 使用路由的方式处理客户端的请求数据](#317-使用路由的方式处理客户端的请求数据)
    - [3.2 使用javaScript实现webSocket客户端](#32-使用javaScript实现webSocket客户端)
        - [3.2.1 启动并注册监听程序](#321-启动并注册监听程序)
        - [3.2.2 发送数据](#322-发送数据)
- [4、goWebSocket 项目](#4goWebSocket-项目)
    - [4.1 项目说明](#41-项目说明)
    - [4.2 项目依赖](#42-项目依赖)
    - [4.3 项目启动](#43-项目启动)
    - [4.4 Nginx配置](#44-Nginx配置)
- [5、webSocket项目Nginx配置](#5webSocket项目Nginx配置)
    - [5.1 nginx配置](#51-nginx配置)
    - [5.2 问题处理](#52-问题处理)
- [6、压测](#6压测)
    - [6.1 Linux内核优化](#61-Linux内核优化)
    - [6.2 压测准备](#62-压测准备)
    - [6.3 压测数据](#63-压测数据)
- [7、如何基于webSocket实现一个分布式Im](#7如何基于webSocket实现一个分布式Im)
    - [7.1 说明](#71-说明)
    - [7.2 架构](#72-架构)
- [8、回顾和反思](#8回顾和反思)
    - [8.1 在其它系统应用](#81-在其它系统应用)
    - [8.2 需要完善、优化](#82-需要完善优化)
    - [8.3 总结](#83-总结)
- [9、参考文献](#9参考文献)


## 1、项目说明
#### 1.1 goWebSocket

golang websocket websocket 中间键，单机支持百万连接，使用gin框架、nginx负载、可以水平部署、程序内部相互通讯、使用grpc通讯协议。

本文将介绍如何实现一个聊天系统(IM)。

#### 1.2 项目体验
- [聊天首页](http://im.91vh.com/home/index)
- 打开连接以后进入聊天界面
- 多人群聊可以同时打开两个窗口

## 2、介绍webSocket
### 2.1 webSocket 是什么
WebSocket 协议在2008年诞生，2011年成为国际标准。所有浏览器都已经支持了。

它的最大特点就是，服务器可以主动向客户端推送信息，客户端也可以主动向服务器发送信息，是真正的双向平等对话，属于服务器推送技术的一种。

![HTTP协议和WebSocket比较](https://img.mukewang.com/5d4cf0750001bc4706280511.png)

![HTTP协议和WebSocket比较](https://img.mukewang.com/5d4cf1180001493404180312.jpg)

- HTTP和webSocket的比较

### 2.2 webSocket的兼容性
- 浏览器的兼容性

![图片描述](https://img.mukewang.com/5d4cf2170001859e12190325.jpg)

- 服务端的支持

golang、java、php、node.js、python、nginx 都有不错的支持

- Android和IOS的支持

Android可以使用java-webSocket对webSocket支持
iOS 4.2及更高版本具有WebSockets支持

### 2.3 为什么要用webSocket
- 从业务上出发
 1. 目前大多数的请求都是使用HTTP，都是由客户端发起一个请求，有服务端处理，然后返回结果，不可以服务端主动向某一个客户端主动发送数据 ![服务端处理一个请求](https://img.mukewang.com/5d4cf5650001773612800720.jpg)
 2. 大多数场景我们需要主动通知用户，如:聊天系统、用户完成任务主动告诉用户、一些运营活动需要通知到在线的用户


## 3、如何实现基于webSocket的长链接系统

### 3.1 使用go实现webSocket服务端

#### 3.1.1 启动端口兼容听
- go实现webSocket可以使用

- **main.go** 实现启动

```
go websocket.StartWebSocket()
```
- **init_acc.go** 启动程序

```
// 启动程序
func StartWebSocket() {

	serverIp = helper.GetServerIp()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort

	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()
	fmt.Println("WebSocket 启动程序成功", serverIp, serverPort)

	http.ListenAndServe(":"+webSocketPort, nil)
}
```

#### 3.1.2 升级协议
- 客户端是通过http请求发送到服务端，我们需要对http协议进行升级
- 对http请求进行升级 [gorilla/websocket](https://github.com/gorilla/websocket) 已经做得很好了，我们直接使用就可以了
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

	conn.CloseHandler()

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
- 当前程序有多少用户连接，还需要对用户广播的需要，这里我们就需要一个管理者，处理这些事件

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
- 防止发生程序崩溃，所以需要捕捉异常
- 为了显示异常崩溃位置这里使用`string(debug.Stack())`打印调用堆栈信息
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
### 3.2 使用javaScript实现webSocket客户端
#### 3.2.1 启动并注册监听程序
- js 建立连接，并处理连接成功、收到数据、断开连接的事件处理

```$js
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
- 建立连接以后由客户端向服务器发送数据示例
```
登录:
ws.send('{"seq":"2323","cmd":"login","data":{"userId":"11","appId":101}}');

心跳:
ws.send('{"seq":"2324","cmd":"heartbeat","data":{}}');
 
关闭连接:
ws.close();
```

## 4、goWebSocket 项目
### 4.1 项目说明
- 本项目是基于webSocket实现的分布式IM系统
- 客户端随机分配用户名，所有人进入一个聊天室，实现群聊的功能
- 单台机器(24核128G内存)支持百万客户端连接
- 支持水平部署，部署的机器之间可以相互通讯

- 项目架构图 (待定)

### 4.2 项目依赖
- 本项目使用govendor管理依赖，克隆本项目就可以直接使用

```
govendor add github.com/gin-gonic/gin@v1.4.0
govendor add -tree github.com/go-redis/redis
govendor add -tree github.com/gorilla/websocket
govendor add -tree github.com/spf13/viper

# viper 依赖
govendor add -tree  github.com/fsnotify/fsnotify
govendor add -tree github.com/hashicorp/hcl
govendor add -tree github.com/magiconair/properties
govendor add -tree github.com/mitchellh/mapstructure
govendor add -tree  github.com/pelletier/go-toml
govendor add -tree  github.com/spf13/afero
govendor add -tree  github.com/spf13/cast
govendor add -tree  github.com/spf13/jwalterweatherman
govendor add -tree  github.com/spf13/pflag
govendor add -tree  github.com/subosito/gotenv
govendor add -tree  golang.org/x/text/transform
govendor add -tree  golang.org/x/text/unicode

# grpc
govendor add -tree google.golang.org/grpc
govendor add -tree google.golang.org/grpc/codes
govendor add -tree github.com/golang/protobuf
govendor add -tree golang.org/x/net
govendor add -tree google.golang.org/genproto/googleapis
govendor add -tree golang.org/x/text/secure/bidirule
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


## 5、webSocket项目Nginx配置
### 5.1 nginx配置
- 使用域名 **im.91vh.com** 为示例，参考配置
- 一级目录**im.91vh.com/acc** 是给webSocket使用，是用nginx流转发功能，转发到golang 8089 端口处理
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


### 5.2 问题处理
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


## 6、压测
### 6.1 Linux内核优化
- 设置文件打开句柄数
```
ulimit -n 1000000
```
- 设置sockets连接参数
```bash
vim /etc/sysctl.conf
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 0
```

### 6.2 压测准备
- 待压测，如果大家有压测的结果欢迎补充

### 6.3 压测数据
| 在线用户数 |   cup  |  内存   |  I/O  | net.out |
| :----:   | :----: | :----: | :----: | :----: |
| 1W       |        |        |        |        |
| 10W      |        |        |        |        |
| 100W     |        |        |        |        |

## 7、如何基于webSocket实现一个分布式Im
### 7.1 说明
- 参考本项目源码，这里不做说明

### 7.2 架构


## 8、回顾和反思
### 8.1 在其它系统应用
- 本系统设计的初衷就是:和客户端保持一个长链接、对外部系统两个接口(查询用户是否在线、给在线的用户推送消息)，实现业务的分离
- 只有和业务分离可，才可以供多个业务使用，而不是每个业务都建立一个长链接，这样消耗的资源有点大

#### 8.2 已经实现的功能

- gin log日志(请求日志+debug日志)
- 读取配置文件 完成
- 定时脚本，清理过期未心跳链接 完成
- http接口，获取登录、链接数量 完成
- http接口，发送push、查询有多少人在线 完成
- grpc 程序内部通讯，发送消息 完成
- appIds 一个用户在多个平台登录
- 界面，把所有在线的人拉倒一个群里面，发送消息 完成
- ~~单聊~~、群聊 完成
- 实现分布式，水平扩张 完成
- 压测脚本
- 文档整理
- 文档目录、百万长链接的实现、为什么要实现一个IM、怎么实现一个Im 
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
- 虽然实现了一个分布式在聊天的IM，但是有很多细节没有处理(登录没有鉴权、界面还待优化等)，但是可以通过这个示例可以了解到通过WebSocket解决很多业务上需求
- 本文虽然号称单台机器能有百万长链接(内存上能满足)，但是实际在场景远比这个复杂(cpu有些压力)，当然了如果你有这么大的业务量可以购买更多的机器更好的去支持你的业务


## 9、参考文献

[维基百科 WebSocket](https://zh.wikipedia.org/wiki/WebSocket)

[阮一峰 WebSocket教程](http://www.ruanyifeng.com/blog/2017/05/websocket.html)

[link1st gowebsocket](https://github.com/link1st/gowebsocket)



