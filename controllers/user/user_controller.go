/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gowebsocket/common"
	"gowebsocket/controllers"
)

// 查看用户是否在线
func Online(c *gin.Context) {

	// 获取参数
	name := c.DefaultPostForm("name", "")
	fmt.Println("Hello", name)

	data := make(map[string]interface{})

	data["name"] = name

	controllers.Response(c, common.OK, "", data)
}

// 给用户发送消息
func SendMessage(c *gin.Context) {

}
