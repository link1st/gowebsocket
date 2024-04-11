// Package models 数据模型
package models

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// Login 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}
