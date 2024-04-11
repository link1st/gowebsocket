// Package models 数据模型
package models

import (
	"errors"
	"fmt"
	"strings"
)

// Server 服务器结构体
type Server struct {
	Ip   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

// NewServer 创建
func NewServer(ip string, port string) *Server {
	return &Server{Ip: ip, Port: port}
}

// String to string
func (s *Server) String() (str string) {
	if s == nil {
		return
	}
	str = fmt.Sprintf("%s:%s", s.Ip, s.Port)
	return
}

// StringToServer 字符串转结构体
func StringToServer(str string) (server *Server, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		return nil, errors.New("err")
	}
	server = &Server{
		Ip:   list[0],
		Port: list[1],
	}
	return
}
