/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 查看用户是否在线
func Index(c *gin.Context) {
	data := gin.H{
		"title": "聊天首页",
	}
	c.HTML(http.StatusOK, "index.tpl", data)
}
