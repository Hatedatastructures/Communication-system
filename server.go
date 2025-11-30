package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip                    string           // Server Ip地址
	Port                  int              // Server 绑定的端口
	MapLock               sync.RWMutex     // Map 锁
	UserMap               map[string]*User // 用户Map
	BroadcastMassageQueue chan string      // Server 接受的的消息队列
}

func NewServer(ip string, port int) *Server {

	NewValue := &Server{
		Ip:                    ip,
		Port:                  port,
		MapLock:               sync.RWMutex{},
		UserMap:               make(map[string]*User),
		BroadcastMassageQueue: make(chan string, 3),
	}
	return NewValue
}

// 监听发送给服务器的消息
func (ThisServer *Server) ListenBroadcast() {
	for {
		BroadcastMassage := <-ThisServer.BroadcastMassageQueue
		ThisServer.BroadcastMassage(BroadcastMassage)
	}
}

// 向所有用户推送消息
func (ThisServer *Server) BroadcastMassage(message string) {
	ThisServer.MapLock.Lock()
	defer ThisServer.MapLock.Unlock()
	for _, UserValue := range ThisServer.UserMap {
		UserValue.MessageQueue <- message
	}
}

// 广播用户消息
func (ThisServer *Server) Broadcast(user *User, message string) {
	SendMessage := fmt.Sprintf("[%v]%v:%v", user.UserName, user.Address, message)
	ThisServer.BroadcastMassageQueue <- SendMessage
}

func (ThisServer *Server) HandleBusiness(conn net.Conn) {
	// 处理业务
	// fmt.Println("新连接...")
	UserValue := NewUser(conn)

	ThisServer.MapLock.Lock()
	ThisServer.UserMap[UserValue.UserName] = UserValue
	ThisServer.MapLock.Unlock()

	ThisServer.Broadcast(UserValue, "欢迎来到服务器")

}

// 启动Server
func (ThisServer *Server) Start() {
	// 拼接完整的IP地址
	var address string = fmt.Sprintf("%s:%d", ThisServer.Ip, ThisServer.Port)
	// 监听端口
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close() // 关闭监听

	go ThisServer.ListenBroadcast() // 监听广播消息
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
