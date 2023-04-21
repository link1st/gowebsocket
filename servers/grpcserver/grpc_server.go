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
	"log"
	"net"

	"github.com/link1st/gowebsocket/common"
	"github.com/link1st/gowebsocket/models"
	"github.com/link1st/gowebsocket/protobuf"
	"github.com/link1st/gowebsocket/servers/websocket"

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

// 查询用户是否在线
func (s *server) QueryUsersOnline(c context.Context,
	req *protobuf.QueryUsersOnlineReq) (rsp *protobuf.QueryUsersOnlineRsp, err error) {

	fmt.Println("grpc_request 查询用户是否在线", req.String())

	rsp = &protobuf.QueryUsersOnlineRsp{}

	online := websocket.CheckUserOnline(req.GetAppId(), req.GetUserId())

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

	data := models.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
	sendResults, err := websocket.SendUserMessageLocal(req.GetAppId(), req.GetUserId(), data)
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

// 给本机全体用户发消息
func (s *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {

	fmt.Println("grpc_request 给本机全体用户发消息", req.String())

	rsp = &protobuf.SendMsgAllRsp{}

	data := models.GetMsgData(req.GetUserId(), req.GetSeq(), req.GetCms(), req.GetMsg())
	websocket.AllSendMessages(req.GetAppId(), req.GetUserId(), data)

	setErr(rsp, common.OK, "")

	fmt.Println("grpc_response 给本机全体用户发消息:", rsp.String())

	return
}

// 获取本机用户列表
func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp,
	err error) {

	fmt.Println("grpc_request 获取本机用户列表", req.String())

	appId := req.GetAppId()
	rsp = &protobuf.GetUserListRsp{}

	// 本机
	userList := websocket.GetUserList(appId)

	setErr(rsp, common.OK, "")
	rsp.UserId = userList

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
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
