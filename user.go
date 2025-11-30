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
}

func NewUser(socket net.Conn) *User {

	UserNameMessage := fmt.Sprint("User-", socket.RemoteAddr().String())

	NewValue := &User{
		UserName:     UserNameMessage,
		Address:      socket.RemoteAddr().String(),
		MessageQueue: make(chan string, 3),
		Connect:      socket,
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
