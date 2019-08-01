/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-01
* Time: 10:40
 */

package models

const (
	messageTypeText = "text"
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
		Type: messageTypeText,
		From: from,
		Msg:  Msg,
	}

	return
}
