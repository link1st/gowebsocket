// Package routers 路由
package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/link1st/gowebsocket/v2/controllers/home"
	"github.com/link1st/gowebsocket/v2/controllers/systems"
	"github.com/link1st/gowebsocket/v2/controllers/user"
)

// Init http 接口路由
func Init(router *gin.Engine) {
	router.LoadHTMLGlob("views/**/*")

	// 用户组
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", user.List)
		userRouter.GET("/online", user.Online)
		userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
	}

	// 系统
	systemRouter := router.Group("/system")
	{
		systemRouter.GET("/state", systems.Status)
	}

	// home
	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
	}
}
