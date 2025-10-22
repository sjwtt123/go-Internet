package ServeMethod

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	same "go-Internet/tcp/Samemethod"
	"io"
	"net"
	"runtime/debug"
	"time"
)

// Client 客户端
type Client struct {
	Conn       net.Conn  // TCP连接
	Nickname   string    // 客户端昵称
	Boo        bool      //是否在线
	LastActive time.Time // 最后活跃时间
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Radio()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	for {
		select {

		case cli := <-leaveChan:
			stringJoin := cli.Nickname + "离开聊天室"
			//加锁
			HbManager.RemoveClient(&cli)
			CloseConn(&cli)
			messageChan <- stringJoin

		case message := <-messageChan:

			fmt.Println(message)
			for _, clientConn := range HbManager.clients {
				err := same.Write(message, clientConn.Conn)
				if err != nil {
					fmt.Printf("客户端%v写入错误：%v\n", clientConn.Nickname, err)
					return
				}
			}

		case privatemessage := <-privateChan:

			err := same.Write(privatemessage.message, privatemessage.client.Conn)
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
				fmt.Printf("与客户端%v连接已关闭\n", client.Nickname)
				leaveChan <- client
				break
			}
			//防止心跳检测已经将连接关闭
			findClient := FindClient(client.Nickname)
			if findClient.Nickname == "" {
				break
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

			message, ty, err := same.AnalyzeMessage(mesa)

			if err != nil {
				return
			}

			switch ty {

			case "UnderLine":
				Underline(client)

			case "Online":
				Online(client)

			case "Private":
				PrivateMessage(message, client)

			case "List":
				SList(client)

			case "Ping":
				//更新心跳检测时间
				HbManager.UpdateClientActivity(client.Nickname)

			default:
				messageChan <- fmt.Sprintf("[%s]:%s", client.Nickname, message)
			}
		}
	}
}

// Online 上线功能
func Online(client Client) {
	HbManager.clients[client.Nickname].Boo = true
	messageChan <- fmt.Sprintf("[%s]已上线", client.Nickname)

}

// Underline 下线功能
func Underline(client Client) {
	HbManager.clients[client.Nickname].Boo = false
	messageChan <- fmt.Sprintf("[%s]已下线", client.Nickname)

}

// SList 列出所有用户
func SList(client Client) {
	var s string
	for c, b := range HbManager.clients {
		fmt.Println(c)
		if b.Boo {
			s += fmt.Sprintf("用户名:%v 状态：已上线\n", c)
		} else {
			s += fmt.Sprintf("用户名:%v 状态：已下线\n", c)
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
	if Sendclient.Nickname == "" {
		privateChan <- ReceiveClient{client: client, message: "未找到该用户无法私发"}
	} else {
		fmt.Printf("用户[%v]对用户[%v]私发了：%v", client.Nickname, name, message)
		privateChan <- ReceiveClient{message: fmt.Sprintf("收到来自[%v]的私发:%v", client.Nickname, message), client: Sendclient}
		privateChan <- ReceiveClient{message: "发送成功", client: client}
	}
}
