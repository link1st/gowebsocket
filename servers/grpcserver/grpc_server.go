// Package grpcserver grpc 服务器
package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/protobuf"
	"github.com/link1st/gowebsocket/v2/servers/websocket"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedAccServerServer
}

func setErr(rsp proto.Message, code uint32, message string) {
	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.QueryUsersOnlineRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgAllRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:
	}
}

// QueryUsersOnline 查询用户是否在线
func (s *server) QueryUsersOnline(c context.Context,
	req *protobuf.QueryUsersOnlineReq) (rsp *protobuf.QueryUsersOnlineRsp, err error) {

	fmt.Println("grpc_request 查询用户是否在线", req.String())

	rsp = &protobuf.QueryUsersOnlineRsp{}

	online := websocket.CheckUserOnline(req.GetAppID(), req.GetUserID())

	setErr(req, common.OK, "")
	rsp.Online = online

	return rsp, nil
}

// SendMsg 给本机用户发消息
func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.SendMsgRsp, err error) {
	fmt.Println("grpc_request 给本机用户发消息", req.String())
	rsp = &protobuf.SendMsgRsp{}
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

// SendMsgAll 给本机全体用户发消息
func (s *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {
	fmt.Println("grpc_request 给本机全体用户发消息", req.String())
	rsp = &protobuf.SendMsgAllRsp{}
	data := models.GetMsgData(req.GetUserID(), req.GetSeq(), req.GetCms(), req.GetMsg())
	websocket.AllSendMessages(req.GetAppID(), req.GetUserID(), data)
	setErr(rsp, common.OK, "")
	fmt.Println("grpc_response 给本机全体用户发消息:", rsp.String())
	return
}

// GetUserList 获取本机用户列表
func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp,
	err error) {

	fmt.Println("grpc_request 获取本机用户列表", req.String())

	appID := req.GetAppID()
	rsp = &protobuf.GetUserListRsp{}

	// 本机
	userList := websocket.GetUserList(appID)

	setErr(rsp, common.OK, "")
	rsp.UserID = userList

	fmt.Println("grpc_response 获取本机用户列表:", rsp.String())

	return
}

// Init rpc server
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
func Init() {
	rpcPort := viper.GetString("app.rpcPort")
	fmt.Println("rpc server 启动", rpcPort)
	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
