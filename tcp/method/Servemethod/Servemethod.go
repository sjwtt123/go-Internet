package Servemethod

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
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

// ClientManager 互斥锁
type ClientManager struct {
	sync.Mutex
	clients map[Client]bool // 在线客户端
}

func (cm *ClientManager) AddClient(cli Client) {
	cm.Lock()
	defer cm.Unlock()
	cm.clients[cli] = true
}

func (cm *ClientManager) RemoveClient(cli Client) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.clients, cli)
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[Client]bool), // 初始化
	}
}

// 使用

// 全局变量
var (
	joinChan    = make(chan Client, 100)        // 客户端加入通道
	leaveChan   = make(chan Client, 100)        // 客户端离开通道
	messageChan = make(chan string, 100)        // 消息广播通道
	privateChan = make(chan ReceiveClient, 100) // 私人消息通道

	clientManager = NewClientManager()
	by            = make([]byte, 1024)
)

// Radio 广播模式，发送信息
func Radio() {

	for {
		select {

		case cli := <-joinChan:
			stringJoin := cli.nickname + "进入聊天室\n"
			fmt.Println(stringJoin)

		case cli := <-leaveChan:
			stringJoin := cli.nickname + "离开聊天室\n"
			//加锁
			clientManager.RemoveClient(cli)

			fmt.Println(stringJoin)
			messageChan <- stringJoin

		case message := <-messageChan:
			by := []byte(message)

			for client := range clientManager.clients {
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

// Createprocess 创建用户信息
func Createprocess(conn net.Conn) error {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭失败")
		}
	}(conn)
LOOP:
	i, err := conn.Read(by)
	if err != nil {
		fmt.Println("读入失败，客户端已关闭")
	}
	scanner := bufio.NewScanner(strings.NewReader(string(by[:i])))
	for scanner.Scan() {
		username := scanner.Text()
		name := IsRepeatName(username)
		if !name {
			_, err2 := conn.Write([]byte("isCreate"))
			if err2 != nil {
				return fmt.Errorf("写入数据失败")
			}
			goto LOOP
		}
		//创建用户对象存每个客户的用户名，和Conn通用的面向流的网络连接
		client := Client{nickname: username, conn: conn}
		clientManager.AddClient(client)
		fmt.Println(clientManager.clients)
		_, err2 := conn.Write([]byte(username + "用户创建成功\n"))
		if err2 != nil {
			return fmt.Errorf("写入数据失败")
		}
		joinChan <- client
		messageChan <- fmt.Sprintf("%s进入聊天室\n", username)
		CRead(conn, client)

	}
	return nil
}

// CRead 客户端读入数据
func CRead(conn net.Conn, client Client) {
	for {
		i, err := conn.Read(by)
		if err != nil {
			fmt.Println("读入失败，客户端已关闭")
			leaveChan <- client
			break
		}
		scanner := bufio.NewScanner(strings.NewReader(string(by[:i])))
		for scanner.Scan() {
			mesa := scanner.Text() // 自动按\n拆分
			if mesa == "" {
				// 处理空消息
				break
			}
			switch {

			case mesa == "EXIT":
				clientManager.clients[client] = false
				messageChan <- fmt.Sprintf("[%s]已下线\n", client.nickname)
				fmt.Printf("[%s]已下线\n", client.nickname)
				continue

			case mesa == "RE":
				clientManager.clients[client] = true
				messageChan <- fmt.Sprintf("[%s]已上线\n", client.nickname)
				fmt.Printf("[%s]已上线\n", client.nickname)
				continue

			case len(mesa) >= 2 && mesa[0:2] == "TO":
				name, message, err1 := FindName(mesa)
				if err1 != nil {
					privateChan <- ReceiveClient{client: client, message: "私发格式错误\n"}
					fmt.Println(client.nickname + "私发输入格式错误")
					continue
				}
				Sendclient := FindClient(name)
				if Sendclient.nickname == "" {
					privateChan <- ReceiveClient{client: client, message: "未找到该用户无法私发\n"}
				} else {
					fmt.Printf("用户[%v]对用户[%v]私发了：%v\n", client.nickname, name, message)
					privateChan <- ReceiveClient{message: fmt.Sprintf("收到来自[%v]的私发:%v\n", client.nickname, message), client: Sendclient}
					privateChan <- ReceiveClient{message: "发送成功\n", client: client}
				}

			case mesa == "LIST":
				SList(client)

			default:
				messageChan <- fmt.Sprintf("[%s]:%s\n", client.nickname, mesa)
				fmt.Printf("[%v]:%s\n", client.nickname, mesa)
			}
		}
	}
}

// IsRepeatName 判断用户名是否重复
func IsRepeatName(username string) bool {
	for client, _ := range clientManager.clients {
		if username == client.nickname {
			return false
		}
	}
	return true
}
func SList(client Client) {
	var s string
	for c, b := range clientManager.clients {
		if b {
			s += fmt.Sprintf("用户名:%v 状态：已上线\n", c.nickname)
		} else {
			s += fmt.Sprintf("用户名:%v 状态：已下线\n", c.nickname)
		}
	}
	receiveClient := ReceiveClient{client: client, message: s}
	privateChan <- receiveClient
}

// FindName 私发功能找用户名
func FindName(s string) (string, string, error) {
	index := strings.Index(s, "[")
	index1 := strings.Index(s, "]")
	if index < 0 || index1 < 0 || index1+2 >= len(s) {
		return "", "", fmt.Errorf("输入格式错误\n")
	}
	message := s[index1+2:]
	s = s[index+1 : index1]

	return s, message, nil
}
func FindClient(name string) Client {
	clientManager.Mutex.Lock()
	defer clientManager.Mutex.Unlock()
	for client, _ := range clientManager.clients {
		if client.nickname == name {
			return client
		}
	}

	fmt.Println("未找到该用户")
	return Client{}

}
