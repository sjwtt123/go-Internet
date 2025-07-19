package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client 客户端结构
type Client struct {
	conn     net.Conn
	nickname string
}

// ReceiveClient 接收客户端用于私发消息,和单个接收消息
type ReceiveClient struct {
	client  Client
	message string
}

// 全局变量
var (
	clients     = make(map[Client]bool)    // 在线客户端
	joinChan    = make(chan Client)        // 客户端加入通道
	leaveChan   = make(chan Client)        // 客户端离开通道
	messageChan = make(chan string)        // 消息广播通道
	privateChan = make(chan ReceiveClient) // 私人消息通道
)

// 广播模式
func radio() {
	for {
		select {

		case cli := <-joinChan:
			stringJoin := cli.nickname + "进入聊天室"
			fmt.Println(stringJoin)

		case cli := <-leaveChan:
			stringJoin := cli.nickname + "离开聊天室\n"
			delete(clients, cli)
			fmt.Println(stringJoin)
			messageChan <- stringJoin

		case message := <-messageChan:
			by := []byte(message)
			for client := range clients {
				_, err := client.conn.Write(by)
				if err != nil {
					fmt.Println("广播写入错误：", err)
				}
			}
		case privatemessage := <-privateChan:
			_, err := privatemessage.client.conn.Write([]byte(privatemessage.message))
			if err != nil {
				fmt.Println("私发消息错误:", err)
			}
		}
	}
}

// 客户端读入数据
func read(conn net.Conn, client Client) {
	for {

		by := make([]byte, 1024)
		read, err := conn.Read(by)

		if err != nil {
			leaveChan <- client
			fmt.Println("远程服务", conn.RemoteAddr().String(), "已退出聊天室")
			return
		}

		s := strings.TrimSpace(string(by)[:read])

		switch {

		case s == "exit":
			clients[client] = false
			messageChan <- fmt.Sprintf("[%s]已下线\n", client.nickname)
			fmt.Printf("[%s]已下线\n", client.nickname)
			continue

		case s == "re":
			clients[client] = true
			messageChan <- fmt.Sprintf("[%s]已上线\n", client.nickname)
			fmt.Printf("[%s]已上线\n", client.nickname)
			continue

		case s[0:2] == "TO":
			name, message, err1 := findName(s)
			if err1 != nil {
				privateChan <- ReceiveClient{client: client, message: "私发格式错误\n"}
				fmt.Println(client.nickname + "私发输入格式错误")
				continue
			}
			Sendclient := findClient(name)
			if Sendclient.nickname == "" {
				privateChan <- ReceiveClient{client: client, message: "未找到该用户无法私发\n"}
			} else {
				fmt.Printf("用户[%v]对用户[%v]私发了：%v\n", client.nickname, name, message)
				privateChan <- ReceiveClient{message: fmt.Sprintf("收到来自[%v]的私发:%v\n", client.nickname, message), client: Sendclient}
				privateChan <- ReceiveClient{message: "发送成功\n", client: client}

			}

		case s == "LIST":
			list(client)

		default:
			messageChan <- fmt.Sprintf("[%s]:%s\n", client.nickname, s)
			fmt.Printf("[%v]:%s\n", client.nickname, s)
		}

	}
}

// Createprocess 创建用户信息
func Createprocess(conn net.Conn) {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭失败")
		}
	}(conn)

LOOP:
	fmt.Println("请创建用户名：")
	readName, err := bufio.NewReader(os.Stdin).ReadString('\n')
	username := strings.TrimSpace(readName)
	name := isRepeatName(username)
	if !name {
		fmt.Println("用户名已存在")
		goto LOOP
	} else if username == "" {
		fmt.Println("用户名不能为空")
		goto LOOP
	}

	if err != nil {
		fmt.Println("输入用户名错误：", err)
	}
	fmt.Println("创建成功")

	//创建用户对象存每个客户的用户名，和Conn通用的面向流的网络连接
	client := Client{nickname: username, conn: conn}
	clients[client] = true

	conn.Write([]byte(username + "用户创建成功\n"))
	joinChan <- client
	messageChan <- fmt.Sprintf("%s进入聊天室\n", username)
	read(conn, client)

}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		fmt.Println("创建服务端失败", err)
	}
	fmt.Println("创建服务端成功，等待连接")

	//延迟关闭
	defer func(listen net.Listener) {
		err1 := listen.Close()
		if err1 != nil {
			fmt.Println("listen关闭失败", err1)
		}
		//返回该接口的网络地址
		fmt.Println(listen.Addr().String())
	}(listen)

	for {
		//等待下一个连接到该接口的连接
		accept, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println("客户端连接失败：", err2)
		}
		fmt.Printf("连接成功对应的ip地址，端口号：%v\n", accept.RemoteAddr().String())

		go radio()

		go Createprocess(accept)

	}
}

// 判断用户名是否重复
func isRepeatName(username string) bool {
	for client, _ := range clients {
		if username == client.nickname {
			return false
		}
	}
	return true
}
func list(client Client) {
	var s string
	for c, b := range clients {
		if b {
			s += fmt.Sprintf("用户名:%v 状态：已上线\n", c.nickname)
		} else {
			s += fmt.Sprintf("用户名:%v 状态：已下线\n", c.nickname)
		}
	}
	receiveClient := ReceiveClient{client: client, message: s}
	privateChan <- receiveClient
}

func findName(s string) (string, string, error) {
	index := strings.Index(s, "[")
	index1 := strings.Index(s, "]")
	if index < 0 || index1 < 0 || index1+2 >= len(s) {
		return "", "", fmt.Errorf("输入格式错误\n")
	}
	message := s[index1+2:]
	s = s[index+1 : index1]

	return s, message, nil
}
func findClient(name string) Client {
	for client, _ := range clients {
		if client.nickname == name {
			return client
		}
	}
	fmt.Println("未找到该用户")
	return Client{}

}
