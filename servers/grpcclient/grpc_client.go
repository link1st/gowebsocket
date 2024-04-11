// Package grpcclient grpc 客户端
package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/protobuf"
)

// SendMsgAll 给全体用户发送消息
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsgAll(server *models.Server, seq string, appID uint32, userID string, cmd string,
	message string) (sendMsgID string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())

		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.SendMsgAllReq{
		Seq:    seq,
		AppID:  appID,
		UserID: userID,
		Cms:    cmd,
		Msg:    message,
	}
	rsp, err := c.SendMsgAll(ctx, &req)
	if err != nil {
		fmt.Println("给全体用户发送消息", err)

		return
	}
	if rsp.GetRetCode() != common.OK {
		fmt.Println("给全体用户发送消息", rsp.String())
		err = errors.New(fmt.Sprintf("发送消息失败 code:%d", rsp.GetRetCode()))

		return
	}
	sendMsgID = rsp.GetSendMsgID()
	fmt.Println("给全体用户发送消息 成功:", sendMsgID)
	return
}

// GetUserList 获取用户列表
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GetUserList(server *models.Server, appID uint32) (userIDs []string, err error) {
	userIDs = make([]string, 0)
	conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.GetUserListReq{
		AppID: appID,
	}
	rsp, err := c.GetUserList(ctx, &req)
	if err != nil {
		fmt.Println("获取用户列表 发送请求错误:", err)
		return
	}
	if rsp.GetRetCode() != common.OK {
		fmt.Println("获取用户列表 返回码错误:", rsp.String())
		err = errors.New(fmt.Sprintf("发送消息失败 code:%d", rsp.GetRetCode()))
		return
	}
	userIDs = rsp.GetUserID()
	fmt.Println("获取用户列表 成功:", userIDs)
	return
}

// SendMsg 发送消息
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsg(server *models.Server, seq string, appID uint32, userID string, cmd string, msgType string,
	message string) (sendMsgID string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.SendMsgReq{
		Seq:     seq,
		AppID:   appID,
		UserID:  userID,
		Cms:     cmd,
		Type:    msgType,
		Msg:     message,
		IsLocal: false,
	}
	rsp, err := c.SendMsg(ctx, &req)
	if err != nil {
		fmt.Println("发送消息", err)
		return
	}
	if rsp.GetRetCode() != common.OK {
		fmt.Println("发送消息", rsp.String())
		err = errors.New(fmt.Sprintf("发送消息失败 code:%d", rsp.GetRetCode()))
		return
	}
	sendMsgID = rsp.GetSendMsgID()
	fmt.Println("发送消息 成功:", sendMsgID)
	return
}
