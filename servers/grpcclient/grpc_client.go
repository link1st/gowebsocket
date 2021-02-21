/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 16:43
 */

package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"gowebsocket/common"
	"gowebsocket/models"
	pb "gowebsocket/protobuf"
	"time"
)

// SendMsg 给用户发送消息
func SendMsg(server *models.Server, seq string, appID uint32, userID string, cmd string, msgType string, message string) (sendMsgID string, err error) {
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("连接失败", server.String())

		return
	}
	defer conn.Close()

	c := pb.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.SendMsgReq{
		Seq:        seq,
		AppID:      appID,
		SendUserID: userID, // 发送者用户ID
		UserID:     userID, // 接收者用户ID
		Cms:        cmd,
		Type:       msgType,
		Msg:        message,
		IsLocal:    false,
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

// rpc client
// 给房间内用户发送消息
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendRoomMsg(server *models.Server, seq string, appID, roomID uint32, userID string, cmd string, message string) (sendMsgID string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("连接失败", server.String())

		return
	}
	defer conn.Close()
	c := pb.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb.SendRoomMsgReq{
		Seq:    seq,
		AppID:  appID,
		RoomID: roomID,
		UserID: userID,
		Cms:    cmd,
		Msg:    message,
	}
	rsp, err := c.SendRoomMsg(ctx, &req)
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

// 获取用户列表
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GetUserList(server *models.Server, appID, roomID uint32) (userIDs []string, err error) {
	userIDs = make([]string, 0)
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer conn.Close()
	c := pb.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb.GetUserListReq{
		AppID:  appID,
		RoomID: roomID,
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
