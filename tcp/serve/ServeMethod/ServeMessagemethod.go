package ServeMethod

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	same "go-Internet/tcp/ReadWritermethod"
	"io"
	"net"
)

// Client 客户端
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
	leaveChan   = make(chan Client, 100)        // 客户端离开通道
	messageChan = make(chan string, 100)        // 消息广播通道
	privateChan = make(chan ReceiveClient, 100) // 私人消息通道

	db *sqlx.DB
)

// Radio 广播模式，发送信息
func Radio() {

	for {
		select {

		case cli := <-leaveChan:
			stringJoin := cli.nickname + "离开聊天室"
			//加锁
			clientManager.RemoveClient(cli)
			CloseConn(cli)
			messageChan <- stringJoin

		case message := <-messageChan:

			fmt.Println(message)
			for client := range clientManager.clients {
				err := same.Write(message, client.conn)
				if err != nil {
					fmt.Printf("客户端%v写入错误：%v\n", client.nickname, err)
					return
				}
			}

		case privatemessage := <-privateChan:

			err := same.Write(privatemessage.message, privatemessage.client.conn)
			if err != nil {
				fmt.Println("私发消息出现错误:", err)
				return
			}
		}
	}
}

// CRead 客户端读入数据处理
func CRead(conn net.Conn, client Client) {
	for {
		scanner, err := same.Read(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("与客户端%v连接已关闭\n", client.nickname)
			}
			leaveChan <- client
			break
		}
		for scanner.Scan() {
			mesa := scanner.Text() // 自动按\n拆分
			if mesa == "" {
				// 处理空消息
				break
			}
			switch {

			case mesa == "EXIT":
				Underline(client)

			case mesa == "RE":
				Online(client)

			case len(mesa) >= 2 && mesa[0:2] == "TO":
				PrivateMessage(mesa, client)

			case mesa == "LIST":
				SList(client)

			default:
				messageChan <- fmt.Sprintf("[%s]:%s", client.nickname, mesa)
			}
		}
	}
}

// SList 列出所有用户
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

// PrivateMessage 处理私发功能
func PrivateMessage(mesa string, client Client) {
	name, message, err1 := FindName(mesa)
	if err1 != nil {
		privateChan <- ReceiveClient{client: client, message: "私发格式错误"}
		return
	}
	Sendclient := FindClient(name)
	if Sendclient.nickname == "" {
		privateChan <- ReceiveClient{client: client, message: "未找到该用户无法私发"}
	} else {
		fmt.Printf("用户[%v]对用户[%v]私发了：%v", client.nickname, name, message)
		privateChan <- ReceiveClient{message: fmt.Sprintf("收到来自[%v]的私发:%v", client.nickname, message), client: Sendclient}
		privateChan <- ReceiveClient{message: "发送成功", client: client}
	}
}

// Online 上线功能
func Online(client Client) {
	clientManager.clients[client] = true
	messageChan <- fmt.Sprintf("[%s]已上线", client.nickname)

}

// Underline 下线功能
func Underline(client Client) {
	clientManager.clients[client] = false
	messageChan <- fmt.Sprintf("[%s]已下线", client.nickname)

}
