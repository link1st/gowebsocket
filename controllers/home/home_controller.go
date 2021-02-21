/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package home

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gowebsocket/helper"
	"gowebsocket/servers/websocket"
	"net/http"
)

// Index 聊天页面
func Index(c *gin.Context) {
	appID := helper.StrToUint32(c.Query("appID"))
	if !websocket.InAppIDs(appID) {
		appID = websocket.GetDefaultAppID()
	}
	fmt.Println("http_request 聊天首页", appID)
	data := gin.H{
		"title":        "聊天首页",
		"appID":        appID,
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.webSocketUrl"),
	}
	c.HTML(http.StatusOK, "index.tpl", data)
}
