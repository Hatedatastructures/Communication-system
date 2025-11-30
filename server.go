package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(Ip string, Port int) *Server {
	return &Server{
		Ip:   Ip,
		Port: Port,
	}
}

func (ThisServer *Server) RegistrationServices() {
	// 注册业务
}

func (ThisServer *Server) HandleBusiness(conn net.Conn) {
	// 处理业务
	fmt.Println("新连接...")
}

func (ThisServer *Server) Start() {
	// 拼接完整的IP地址
	var address string = fmt.Sprintf("%s:%d", ThisServer.Ip, ThisServer.Port)
	// Listen
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close() // 关闭监听
	// 获取套接字
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err:", err)
			continue
		}
		go ThisServer.HandleBusiness(conn)
	}
	// 业务代码
}
