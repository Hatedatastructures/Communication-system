package main

import (
	"fmt"
	"net"
)

type User struct {
	UserName     string      // 名称
	Address      string      // IP 地址
	MessageQueue chan string // 消息队列
	Connect      net.Conn    // 连接
	Server       *Server     // 服务器
}

func NewUser(socket net.Conn, server *Server) *User {

	UserNameMessage := fmt.Sprint("User-", socket.RemoteAddr().String())

	NewValue := &User{
		UserName:     UserNameMessage,
		Address:      socket.RemoteAddr().String(),
		MessageQueue: make(chan string, 3),
		Connect:      socket,
		Server:       server,
	}

	go NewValue.ListenMessage()

	return NewValue
}

func (ThisUser *User) ListenMessage() {
	for {
		BroadcastMassage := <-ThisUser.MessageQueue
		ThisUser.Connect.Write([]byte(BroadcastMassage + "\n"))
	}
}

// 用户上线
func (ThisUser *User) GoOnline() {
	ThisUser.Server.MapLock.Lock()
	ThisUser.Server.UserMap[ThisUser.UserName] = ThisUser
	ThisUser.Server.MapLock.Unlock()
	ThisUser.Server.Broadcast(ThisUser, ThisUser.UserName+"上线了")
}

// 用户下线
func (ThisUser *User) Offline() {
	ThisUser.Server.Broadcast(ThisUser, ThisUser.UserName+"下线了")
	ThisUser.Server.MapLock.Lock()
	delete(ThisUser.Server.UserMap, ThisUser.UserName)
	ThisUser.Server.MapLock.Unlock()
	ThisUser.Connect.Close()
}

// 发送消息
func (ThisUser *User) SendMessage(message string) {
	ThisUser.Server.Broadcast(ThisUser, message)
}
