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
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gowebsocket/common"
	"gowebsocket/protobuf"
	"gowebsocket/servers/users"
	"log"
	"net"
)

type server struct {
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
	case *protobuf.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

// 查询用户是否在线
func (s *server) QueryUsersOnline(c context.Context, req *protobuf.QueryUsersOnlineReq) (rsp *protobuf.QueryUsersOnlineRsp, err error) {

	fmt.Println("grpc_request 查询用户是否在线", req.String())

	rsp = &protobuf.QueryUsersOnlineRsp{}

	online := users.CheckUserOnline(req.GetAppId(), req.GetUserId())

	setErr(req, common.OK, "")
	rsp.Online = online

	return rsp, nil
}

// 给本机用户发消息
func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.SendMsgRsp, err error) {

	fmt.Println("grpc_request 给本机用户发消息", req.String())

	rsp = &protobuf.SendMsgRsp{}

	if req.GetIsLocal() {

		// 不支持
		setErr(rsp, common.ParameterIllegal, "")

		return
	}

	sendResults, err := users.SendUserMessage(req.GetAppId(), req.GetUserId(), req.GetMsg(), req.GetMsg())
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

	return
}

// 给本机用户发消息
func (s *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {

	fmt.Println("grpc_request 给本机用户发消息", req.String())

	rsp = &protobuf.SendMsgAllRsp{}

	sendResults, err := users.SendUserMessageAll(req.GetAppId(), req.GetUserId(), req.GetMsg(), req.GetMsg())
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

	return
}

// 获取用户列表
func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp, err error) {

	fmt.Println("grpc_request 获取用户列表", req.String())

	rsp = &protobuf.GetUserListRsp{}

	userList := users.UserList()

	setErr(rsp, common.OK, "")
	rsp.UserId = userList

	return nil, nil
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
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
