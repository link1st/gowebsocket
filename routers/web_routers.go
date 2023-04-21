/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/link1st/gowebsocket/controllers/home"
	"github.com/link1st/gowebsocket/controllers/systems"
	"github.com/link1st/gowebsocket/controllers/user"
)

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

	// router.POST("/user/online", user.Online)
}
