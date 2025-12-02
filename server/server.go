package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip                    string           // Server Ip地址
	Port                  int              // Server 绑定的端口
	MapLock               sync.RWMutex     // Map 锁
	UserMap               map[string]*User // 用户Map
	BroadcastMassageQueue chan string      // Server 接受的的消息队列
}

// NewServer 创建服务器
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

// ListenBroadcast 监听发送给服务器的消息
func (ThisServer *Server) ListenBroadcast() {
	for {
		BroadcastMassage := <-ThisServer.BroadcastMassageQueue
		ThisServer.BroadcastMassage(BroadcastMassage)
	}
}

// BroadcastMassage 向所有用户推送消息
func (ThisServer *Server) BroadcastMassage(message string) {
	ThisServer.MapLock.Lock()
	defer ThisServer.MapLock.Unlock()
	for _, UserValue := range ThisServer.UserMap {
		UserValue.MessageQueue <- message
	}
}

// Broadcast 广播用户消息
func (ThisServer *Server) Broadcast(user *User, message string) {
	SendMessage := fmt.Sprintf("[%v]%v:%v", user.Address, user.UserName, message)
	ThisServer.BroadcastMassageQueue <- SendMessage
}

// MonitorUserMessages 监听用户的消息
func (ThisServer *Server) MonitorUserMessages(user *User, Condition chan bool) {
	for {
		UserMessage := make([]byte, 4096)
		n, err := user.Connect.Read(UserMessage)

		if n == 0 { // 用户断开连接
			fmt.Println(user.UserName, "断开连接")
			user.Offline() // 下线
			return         // 退出循环防止无限打印
		}

		if err != nil {
			fmt.Println("Read err:", err)
			continue
		}
		message := string(UserMessage[:n-1])
		user.DoMessage(message)
		Condition <- true
	}
}

// HandleBusiness 业务处理
func (ThisServer *Server) HandleBusiness(conn net.Conn) {
	// 处理业务
	// fmt.Println("新连接...")
	UserValue := NewUser(conn, ThisServer)
	UserValue.GoOnline() // 上线

	TimeoutCondition := make(chan bool)

	go ThisServer.MonitorUserMessages(UserValue, TimeoutCondition)

	for {
		select {
		case <-TimeoutCondition:

		case <-time.After(time.Second * 600):
			UserValue.SendNetworkMessage("超时下线")
			close(UserValue.MessageQueue) //销毁资源
			conn.Close()
			return
		}
	}
}

// Start 启动Server
func (ThisServer *Server) Start() {
	// 拼接完整的IP地址
	var address = fmt.Sprintf("%s:%d", ThisServer.Ip, ThisServer.Port)
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
