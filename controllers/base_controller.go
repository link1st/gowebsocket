/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package controllers

import (
	"github.com/gin-gonic/gin"
	"gowebsocket/common"
	"net/http"
)

type BaseController struct {
	gin.Context
}

// 获取全部请求解析到map
func Response(c *gin.Context, code int, msg string, data map[string]interface{}) {
	message := common.Response(code, msg, data)
	c.JSON(http.StatusOK, message)

	return
}
