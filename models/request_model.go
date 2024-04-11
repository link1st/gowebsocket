// Package models 数据模型
package models

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一ID
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// Login 登录请求数据
type Login struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppID        uint32 `json:"appID,omitempty"`
	UserID       string `json:"userID,omitempty"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserID string `json:"userID,omitempty"`
}
