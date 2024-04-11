// Package main 实现一个基于websocket分布式聊天(IM)系统。
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/link1st/gowebsocket/v2/lib/redislib"
	"github.com/link1st/gowebsocket/v2/routers"
	"github.com/link1st/gowebsocket/v2/servers/grpcserver"
	"github.com/link1st/gowebsocket/v2/servers/task"
	"github.com/link1st/gowebsocket/v2/servers/websocket"
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
	_ = http.ListenAndServe(":"+httpPort, router)
}

// initFile 初始化日志
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
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}

func initRedis() {
	redislib.NewClient()
}

func open() {
	time.Sleep(1000 * time.Millisecond)
	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"
	fmt.Println("访问页面体验:", httpUrl)
	cmd := exec.Command("open", httpUrl)
	_, _ = cmd.Output()
}
