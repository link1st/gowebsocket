## 目录
- [goWebSocket](#goWebSocket)
- [架构图](#架构图)
 * 二级目录

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
```

#### 待办事项
- 读取配置文件
- 定时脚本，清理过期未心跳链接
- http接口，获取登录、链接数量
- grpc 程序内部通讯，发送消息
- http接口，发送push、查询有多少人在线
- appIds 一个用户在多个平台登录
- 界面，把所有在线的人拉倒一个群里面，发送消息
- 群聊、单聊