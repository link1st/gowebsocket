/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 09:59
 */

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gowebsocket/lib/redislib"
	"gowebsocket/routers"
	"gowebsocket/servers/grpcserver"
	"gowebsocket/servers/task"
	"gowebsocket/servers/websocket"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	initConfig()

	initFile()

	initRedis()

	router := gin.Default()
	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()

	// 定时任务
	task.Init()

	// 服务注册
	task.ServerInit()

	go websocket.StartWebSocket()
	// grpc
	go grpcserver.Init()

	go open()

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)

}

// 初始化日志
func initFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // 添加搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}

func initRedis() {
	redislib.ExampleNewClient()
}

func open() {

	time.Sleep(1000 * time.Millisecond)

	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"

	fmt.Println("访问页面体验:", httpUrl)

	cmd := exec.Command("open", httpUrl)
	cmd.Output()
}
