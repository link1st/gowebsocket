/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:20
 */

package routers

import (
	"github.com/gin-gonic/gin"
	"gowebsocket/controllers/user"
)

func Init(router *gin.Engine) {

	// 用户组
	userRouter := router.Group("/user")
	{
		userRouter.POST("/online", user.Online)
		userRouter.POST("/sendMessage", user.SendMessage)
	}

	// router.POST("/user/online", user.Online)
}
