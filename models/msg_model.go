/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-01
* Time: 10:40
 */

package models

import "github.com/link1st/gowebsocket/common"

const (
	MessageTypeText = "text"

	MessageCmdMsg   = "msg"
	MessageCmdEnter = "enter"
	MessageCmdExit  = "exit"
)

// 消息的定义
type Message struct {
	Target string `json:"target"` // 目标
	Type   string `json:"type"`   // 消息类型 text/img/
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
}

func NewTestMsg(from string, Msg string) (message *Message) {

	message = &Message{
		Type: MessageTypeText,
		From: from,
		Msg:  Msg,
	}

	return
}

func getTextMsgData(cmd, uuId, msgId, message string) string {
	textMsg := NewTestMsg(uuId, message)
	head := NewResponseHead(msgId, cmd, common.OK, "Ok", textMsg)

	return head.String()
}

// 文本消息
func GetMsgData(uuId, msgId, cmd, message string) string {

	return getTextMsgData(cmd, uuId, msgId, message)
}

// 文本消息
func GetTextMsgData(uuId, msgId, message string) string {

	return getTextMsgData("msg", uuId, msgId, message)
}

// 用户进入消息
func GetTextMsgDataEnter(uuId, msgId, message string) string {

	return getTextMsgData("enter", uuId, msgId, message)
}

// 用户退出消息
func GetTextMsgDataExit(uuId, msgId, message string) string {

	return getTextMsgData("exit", uuId, msgId, message)
}
