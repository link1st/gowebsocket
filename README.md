## 目录
- [goWebSocket](#goWebSocket)
- [架构图](#架构图)
 * 二级目录

#### 体验地址
- [聊天首页](http://im.91vh.com/home/index)

### goWebSocket
golang websocket websocket 中间键，单机支持百万连接，使用gin框架、nginx负载、可以水平部署、程序内部相互通讯、使用grpc通讯协议

### 架构图


### QA


一个用户，多个平台同时登陆
发消息的时候选择一个平台发送


### js发送消息
```$xslt
ws = new WebSocket("ws://127.0.0.1:8089/acc");
 
// setTimeout(时间，"JS代码");
 
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

### 发送数据
```$xslt
登录:
ws.send('{"seq":"2323","cmd":"login","data":{"userId":"11","appId":101}}');

心跳:
ws.send('{"seq":"2324","cmd":"heartbeat","data":{}}');
 
关闭连接:
ws.close();
```

#### goVendor
```bash
govendor add github.com/gin-gonic/gin@v1.4.0
govendor add -tree github.com/go-redis/redis
govendor add -tree github.com/gorilla/websocket
govendor add -tree github.com/spf13/viper
# 依赖
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
govendor add -tree  golang.org/x/text
govendor add -tree  golang.org/x/text/unicode
```


#### 部署
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
 # 修改项目启动地址，redis连接等 如果不需要修改可以跳过
 vim app.yaml
 # 返回项目目录，为以后启动做准备
 cd ..
```

- 启动项目
```
go run main.go
```
- 设置域名、或者是host

- im 域名配置
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

- nginx 错误
```
/link/server/tengine/sbin/nginx -t
nginx: [emerg] unknown "connection_upgrade" variable
configuration file /link/server/tengine/conf/nginx.conf test failed
```

- 需要在 nginx.conf 添加下面代码

```
http{
	fastcgi_temp_file_write_size 128k;
.....

    #support websocket
    map $http_upgrade $connection_upgrade {
        default upgrade;
        ''      close;
    }

.....
    gzip on;
    
}


```

#### 待办事项
- gin log日志(请求日志+debug日志)
- 读取配置文件 完成
- 定时脚本，清理过期未心跳链接 完成
- http接口，获取登录、链接数量 完成
- http接口，发送push、查询有多少人在线 完成
- grpc 程序内部通讯，发送消息
- appIds 一个用户在多个平台登录
- 界面，把所有在线的人拉倒一个群里面，发送消息 完成
- ~~单聊~~、群聊
- 实现分布式，水平扩张

#### 小项
- 定义文本消息结构 完成
- html发送文本消息 完成
- 接口接收文本消息并发送给全体 完成
- html接收到消息 显示到界面 完成
- 界面优化 需要持续优化
- 有人加入以后广播全体 完成
- 定义加入聊天室的消息结构 完成
- 引入机器人 待定