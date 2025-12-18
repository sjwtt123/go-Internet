package ClientMethod

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
	"go-Internet/tcp/tool/redis"
	"log"
	"net"
	"runtime/debug"
	"time"
)

type Client struct {
	Boo      bool //是否上线
	Dial     net.Conn
	Nickname string //用户名
}

// Read 并发读协程
func (client *Client) Read() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Read()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	for {
		msg, err := same.Read(client.Dial)
		if err != nil {
			return fmt.Errorf("读取结束,连接关闭:%v", err)
		}
		fmt.Println(msg)
	}
}

// WriteHeart 并发写协程来进行心跳检测
func (client *Client) WriteHeart() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("WriteHeart()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		message, err := same.CreateMessage(client.Nickname, "", TypeHeart, "")
		if err != nil {
			return fmt.Errorf("心跳检测ping客户端:%v", err)
		}

		err = same.Write(message, client.Dial)
		if err != nil {
			return fmt.Errorf("心跳检测ping客户端:%v", err)
		}
	}
	return nil

}

// Start 开启业务
func (client *Client) Start() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf(" Start()协程发生panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	//读取离线消息
	err1 := redis.ReceiveHistoryAllMessage(client.Nickname)
	if err1 != nil {
		return err1
	}
	fmt.Println("------以上是历史消息------")

	//欢迎开场白
	showWelcomeMessage()
	for {
		//从终端取到数据
		readString, err := GetfromShell()
		if err != nil {
			return err
		}

		switch readString {
		case "1":
			showPrivateChatHelp()

		case "2":
			client.Boo, err = client.Online()
			if err != nil {
				return err
			}
			continue

		case "3":
			client.Boo, err = client.UnderLine()
			if err != nil {
				return err
			}
			continue

		case "4":
			err1 := client.List()
			if err1 != nil {
				return err
			}

		case "5":
			err = ActiveList()
			if err != nil {
				return err
			}

		case "6":
			err := client.Leave()
			if err != nil {
				return fmt.Errorf("关闭与服务端连接失败:%v", err)
			}
			fmt.Println("已退出聊天室")
			return nil

		default:
			err2 := client.WriteTO(readString)
			if err2 != nil {
				return err
			}
		}
	}
}

// Online 上线功能
func (client *Client) Online() (bool, error) {
	if client.Boo {
		fmt.Println("已上线，请勿重复上线功能")
	} else {
		fmt.Println("已上线，可以与别人交流")

		marshal, err2 := same.CreateMessage(client.Nickname, "", TypeOnline, "")
		if err2 != nil {
			return false, fmt.Errorf("客户端创建上线功能：%v", err2)
		}

		err := same.Write(marshal, client.Dial)
		if err != nil {
			return false, fmt.Errorf("客户端创建上线功能：%v", err)
		}

		return true, nil
	}
	return true, nil
}

// UnderLine 下线功能
func (client *Client) UnderLine() (bool, error) {
	if !client.Boo {
		fmt.Println("已下线，请勿重复上线功能")
	} else {
		fmt.Println("已下线，无法回复消息")

		marshal, err2 := same.CreateMessage(client.Nickname, "", TypeUnderLine, "")
		if err2 != nil {
			return false, fmt.Errorf("客户端创建下线功能：%v", err2)
		}

		//将数据写到流内
		err := same.Write(marshal, client.Dial)
		if err != nil {
			return false, fmt.Errorf("客户端创建下线功能：%v", err)
		}
		client.Boo = false
		return client.Boo, nil
	}
	return false, nil
}

// WriteTO 发送功能
func (client *Client) WriteTO(readString string) error {
	if !client.Boo {
		fmt.Println("正在下线状态无法发送信息")
		return nil
	}
	var marshal string
	var err1 error

	//判断是私发还是群发功能(不应该根据前面两个字符来判断私发，因该有流程选择开启私发过程)
	if len(readString) >= 2 && readString[0:2] == "TO" {

		marshal, err1 = same.CreateMessage(client.Nickname, "", TypePrivate, readString)
		if err1 != nil {
			return fmt.Errorf("客户端私发消息功能：%v", err1)
		}

	} else {
		marshal, err1 = same.CreateMessage(client.Nickname, "", TypeRadio, readString)
		if err1 != nil {
			return fmt.Errorf("客户端群发消息功能：%v", err1)
		}
	}

	err := same.Write(marshal, client.Dial)
	if err != nil {
		return fmt.Errorf("客户端发送消息功能：%v", err)
	}

	return nil
}

// List 列表功能
func (client *Client) List() error {

	marshal, err1 := same.CreateMessage(client.Nickname, "", TypeList, "")
	if err1 != nil {
		return fmt.Errorf("客户端list列表功能：%v", err1)
	}

	err := same.Write(marshal, client.Dial)
	if err != nil {
		return fmt.Errorf("客户端list列表功能：%v", err)
	}
	return nil
}

// Leave 列表功能
func (client *Client) Leave() error {

	marshal, err1 := same.CreateMessage(client.Nickname, "", "leave", "")
	if err1 != nil {
		return fmt.Errorf("客户端退出聊天室功能：%v", err1)
	}

	err := same.Write(marshal, client.Dial)
	if err != nil {
		return fmt.Errorf("客户端退出聊天室功能：%v", err)
	}
	return nil
}

func ActiveList() error {
	active, err := redis.FindAllActive()
	if err != nil {
		return err
	}
	fmt.Println("----------------------------")
	for i, a := range active {
		fmt.Printf("活跃度排名：%v  用户：%v\n", i+1, a.Member.(string))
	}
	fmt.Println("----------------------------")
	return nil
}
