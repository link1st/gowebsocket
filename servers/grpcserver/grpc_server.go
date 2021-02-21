/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 16:43
 */

package grpcserver

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gowebsocket/common"
	"gowebsocket/models"
	pb "gowebsocket/protobuf"
	"gowebsocket/servers/websocket"
	"log"
	"net"
)

// server server
type server struct {
	pb.UnimplementedAccServerServer
}

// setErr 设置错误信息
func setErr(rsp proto.Message, code uint32, message string) {
	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *pb.QueryUsersOnlineRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *pb.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *pb.SendRoomMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *pb.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:
	}
}

// 查询用户是否在线
func (s *server) QueryUsersOnline(_ context.Context, req *pb.QueryUsersOnlineReq) (rsp *pb.QueryUsersOnlineRsp, err error) {
	searchRspStr, _ := (&jsonpb.Marshaler{}).MarshalToString(req)
	fmt.Println("grpc_request 查询用户是否在线", searchRspStr)

	// 发送请求
	rsp = &pb.QueryUsersOnlineRsp{}
	online := websocket.CheckUserOnline(req.GetAppID(), req.GetUserID())
	setErr(req, common.OK, "")
	rsp.Online = online
	return rsp, nil
}

// 给本机用户发消息
func (s *server) SendMsg(_ context.Context, req *pb.SendMsgReq) (rsp *pb.SendMsgRsp, err error) {
	fmt.Println("grpc_request 给本机用户发消息", req.String())
	rsp = &pb.SendMsgRsp{}
	if req.GetIsLocal() {
		// 不支持
		setErr(rsp, common.ParameterIllegal, "")
		return
	}

	data := models.GetMsgData(req.GetUserID(), req.GetSeq(), req.GetCms(), req.GetMsg())
	sendResults, err := websocket.SendUserMessageLocal(req.GetAppID(), req.GetUserID(), data)
	if err != nil {
		fmt.Println("系统错误", err)
		setErr(rsp, common.ServerError, "")
		return rsp, nil
	}
	if !sendResults {
		fmt.Println("发送失败", err)
		setErr(rsp, common.OperationFailure, "")
		return rsp, nil
	}
	setErr(rsp, common.OK, "")
	fmt.Println("grpc_response 给本机用户发消息", rsp.String())
	return
}

// SendRoomMsg 给房间内用法发送消息
func (s *server) SendRoomMsg(_ context.Context, req *pb.SendRoomMsgReq) (rsp *pb.SendRoomMsgRsp, err error) {
	fmt.Println("grpc_request 给房间内用法发送消息", req.String())
	rsp = &pb.SendRoomMsgRsp{}
	data := models.GetMsgData(req.GetUserID(), req.GetSeq(), req.GetCms(), req.GetMsg())
	websocket.AllSendMessages(req.GetAppID(), req.GetUserID(), data)
	setErr(rsp, common.OK, "")
	fmt.Println("grpc_response 给房间内用法发送消息:", rsp.String())

	return
}

// 获取本机用户列表
func (s *server) GetUserList(_ context.Context, req *pb.GetUserListReq) (rsp *pb.GetUserListRsp, err error) {
	fmt.Println("grpc_request 获取本机用户列表", req.String())
	appID := req.GetAppID()
	rsp = &pb.GetUserListRsp{}
	// 本机
	userList := websocket.GetUserList(appID)
	setErr(rsp, common.OK, "")
	rsp.UserID = userList
	fmt.Println("grpc_response 获取本机用户列表:", rsp.String())
	return
}

// rpc server
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
func Init() {
	rpcPort := viper.GetString("app.rpcPort")
	fmt.Println("rpc server 启动", rpcPort)
	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
