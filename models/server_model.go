/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 15:38
 */

package models

import (
	"errors"
	"fmt"
	"strings"
)

type Server struct {
	Ip   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

func NewServer(ip string, port string) *Server {

	return &Server{Ip: ip, Port: port}
}

func (s *Server) String() (str string) {
	if s == nil {
		return
	}

	str = fmt.Sprintf("%s:%s", s.Ip, s.Port)

	return
}

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
