package main

import (
	"fmt"
	"net"
	"strings"
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
		_, err := ThisUser.Connect.Write([]byte(BroadcastMassage + "\n"))
		if err != nil {
			return
		}
	}
}

// GoOnline 用户上线
func (ThisUser *User) GoOnline() {
	ThisUser.Server.MapLock.Lock()
	ThisUser.Server.UserMap[ThisUser.UserName] = ThisUser
	ThisUser.Server.MapLock.Unlock()
	ThisUser.Server.Broadcast(ThisUser, ThisUser.UserName+"上线了")
}

// Offline 用户下线
func (ThisUser *User) Offline() {
	ThisUser.Server.Broadcast(ThisUser, ThisUser.UserName+"下线了")
	ThisUser.Server.MapLock.Lock()
	delete(ThisUser.Server.UserMap, ThisUser.UserName)
	ThisUser.Server.MapLock.Unlock()
	err := ThisUser.Connect.Close()
	if err != nil {
		return
	}
}

func (ThisUser *User) SendMessage(message string) {
	ThisUser.MessageQueue <- message
}

func (ThisUser *User) SendNetworkMessage(message string) {
	ThisUser.Connect.Write([]byte(message + "\n"))
}

// DoMessage 发送消息
func (ThisUser *User) DoMessage(message string) {
	if message == "who" { // 查询
		ThisUser.SendMessage("在线用户列表：")
		for UserName := range ThisUser.Server.UserMap {
			UserMessage := fmt.Sprintf("[%v]%v 在线", UserName, ThisUser.Address)
			ThisUser.SendMessage(UserMessage)
		}
	} else if len(message) > 7 && message[:7] == "rename|" { // 修改用户名
		NewUserName := strings.Split(message, "|")[1]
		_, ok := ThisUser.Server.UserMap[NewUserName]
		if ok {
			ThisUser.SendMessage("当前用户名已存在")
		} else {
			ThisUser.Server.MapLock.Lock()
			delete(ThisUser.Server.UserMap, ThisUser.UserName)
			ThisUser.UserName = NewUserName
			ThisUser.Server.UserMap[ThisUser.UserName] = ThisUser
			ThisUser.Server.MapLock.Unlock()
			ThisUser.SendMessage("用户名已更改为：" + NewUserName)
		}
	} else {
		ThisUser.Server.Broadcast(ThisUser, message)
	}
}
