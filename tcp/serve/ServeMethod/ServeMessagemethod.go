package ServeMethod

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	same "go-Internet/tcp/Samemethod"
	"go-Internet/tcp/tool/redis"
	"log"
	"net"
	"runtime/debug"
	"time"
)

// Client 客户端
type Client struct {
	Conn       net.Conn // TCP连接
	Nickname   string   // 客户端昵称
	Boo        bool     //是否在线
	LastActive time.Time
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
func Radio(ctx context.Context) {
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
					log.Printf("客户端%v写入错误：%v\n", clientConn.Nickname, err)
					return
				}
			}

		case privatemessage := <-privateChan:

			err := same.Write(privatemessage.message, privatemessage.client.Conn)
			if err != nil {
				log.Println("私发消息出现错误:", err)
				return
			}
		case <-ctx.Done():
			log.Println("Radio 协程退出")
			return

		}
	}
}

// CRead 客户端读入数据处理
func (client *Client) CRead() error {
	for {
		//接受实时消息
		realtimeMessage, err := redis.ReceiveRealtimeMessage(client.Nickname)
		if err != nil {
			return err
		}
		if realtimeMessage == "" {
			// 处理空消息
			continue
		}
		//解析消息
		message, ty, err := same.AnalyzeMessage(realtimeMessage)
		if err != nil {
			return err
		}

		//发送一条消息增加一次活跃度
		err = redis.IncreaseActive(client.Nickname)
		if err != nil {
			return err
		}

		switch ty {

		case TypeUnderLine:
			client.Underline()

		case TypeOnline:
			client.Online()

		case TypePrivate:
			client.PrivateMessage(message)

		case TypeList:
			client.SList()

		case TypeHeart:
			//更新心跳检测时间
			HbManager.UpdateClientActivity(client.Nickname)

		case TypeRadio:
			//群发功能存到list中
			marsh, err1 := same.CreateMessage(client.Nickname, "", TypeRadio, message)
			if err != nil {
				return err
			}
			err1 = redis.AddHistoryToList(marsh)
			if err1 != nil {
				return err
			}
			messageChan <- fmt.Sprintf("[%s]:%s", client.Nickname, message)

		default:
			fmt.Printf("与客户端%v连接已关闭\n", client.Nickname)
			leaveChan <- *client
			return fmt.Errorf("程序结束提示错误")
		}
	}

}

// Online 上线功能
func (client *Client) Online() {
	HbManager.clients[client.Nickname].Boo = true
	messageChan <- fmt.Sprintf("[%s]已上线", client.Nickname)

}

// Underline 下线功能
func (client *Client) Underline() {
	HbManager.clients[client.Nickname].Boo = false
	messageChan <- fmt.Sprintf("[%s]已下线", client.Nickname)

}

// SList 列出所有用户
func (client *Client) SList() {
	var s string
	for c, b := range HbManager.clients {
		fmt.Println(c)
		if b.Boo {
			s += fmt.Sprintf("用户名:%v 状态：已上线\n", c)
		} else {
			s += fmt.Sprintf("用户名:%v 状态：已下线\n", c)
		}
	}
	receiveClient := ReceiveClient{client: *client, message: s}
	privateChan <- receiveClient
}

// PrivateMessage 处理私发功能
func (client *Client) PrivateMessage(mesa string) {

	name, message, err1 := FindName(mesa)
	if err1 != nil {
		privateChan <- ReceiveClient{client: *client, message: "私发格式错误"}
		return
	}

	Sendclient := FindClient(name)
	if Sendclient.Nickname == "" {
		privateChan <- ReceiveClient{client: *client, message: "未找到该用户无法私发"}
	} else {

		//私发功能存到历史消息list中
		marsh, err2 := same.CreateMessage(client.Nickname, Sendclient.Nickname, TypeRadio, message)
		if err2 != nil {
			return
		}
		err1 = redis.AddHistoryToList(marsh)
		if err1 != nil {
			return
		}

		fmt.Printf("用户[%v]对用户[%v]私发了：%v\n", client.Nickname, name, message)
		privateChan <- ReceiveClient{message: fmt.Sprintf("收到来自[%v]的私发:%v", client.Nickname, message), client: Sendclient}
		privateChan <- ReceiveClient{message: "发送成功", client: *client}
	}
}
