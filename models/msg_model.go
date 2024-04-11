// Package models 数据模型
package models

import "github.com/link1st/gowebsocket/v2/common"

const (
	// MessageTypeText 文本类型消息
	MessageTypeText = "text"
	// MessageCmdMsg 文本类型消息
	MessageCmdMsg = "msg"
	// MessageCmdEnter 用户进入类型消息
	MessageCmdEnter = "enter"
	// MessageCmdExit 用户退出类型消息
	MessageCmdExit = "exit"
)

// Message 消息的定义
type Message struct {
	Target string `json:"target"` // 目标
	Type   string `json:"type"`   // 消息类型 text/img/
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
}

// NewMsg 创建新的消息
func NewMsg(from string, Msg string) (message *Message) {
	message = &Message{
		Type: MessageTypeText,
		From: from,
		Msg:  Msg,
	}
	return
}

func getTextMsgData(cmd, uuID, msgID, message string) string {
	textMsg := NewMsg(uuID, message)
	head := NewResponseHead(msgID, cmd, common.OK, "Ok", textMsg)

	return head.String()
}

// GetMsgData 文本消息
func GetMsgData(uuID, msgID, cmd, message string) string {
	return getTextMsgData(cmd, uuID, msgID, message)
}

// GetTextMsgData 文本消息
func GetTextMsgData(uuID, msgID, message string) string {
	return getTextMsgData("msg", uuID, msgID, message)
}

// GetTextMsgDataEnter 用户进入消息
func GetTextMsgDataEnter(uuID, msgID, message string) string {
	return getTextMsgData("enter", uuID, msgID, message)
}

// GetTextMsgDataExit 用户退出消息
func GetTextMsgDataExit(uuID, msgID, message string) string {
	return getTextMsgData("exit", uuID, msgID, message)
}
