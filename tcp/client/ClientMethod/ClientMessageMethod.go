package ClientMethod

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
	"net"
	"runtime/debug"
	"time"
)

var (
	boo  = true
	dial net.Conn
)

// Read 并发读协程
func Read(dial net.Conn) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Read()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	for {
		scanner, err := same.Read(dial)
		if err != nil {
			fmt.Println("读取结束,连接关闭")
			return err
		}

		// 检查退出原因：错误或EOF
		if err = scanner.Err(); err != nil {
			fmt.Println("连接关闭，读取结束")
			return err
		}
		for scanner.Scan() {
			msg := scanner.Text() // 自动按\n拆分
			fmt.Println(msg)
		}

	}
}

// WriteHeart 并发写协程来进行心跳检测
func WriteHeart() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("WriteHeart()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		message, err := same.CreateMessage("Ping", "")
		if err != nil {
			fmt.Println("心跳检测ping客户端创建消息失败", err)
			return err
		}

		err1 := same.Write(message, dial)
		if err1 != nil {
			fmt.Println("心跳检测客户端ping服务端传入数据失败：", err1)
			return err1
		}
	}
	return nil

}

// Start 开启业务
func Start() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(" Start()协程发生panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	//欢迎开场白
	showWelcomeMessage()

	for {
		//从终端取到数据
		readString, err := GetfromShell()
		if err != nil {
			fmt.Println("终端输入失败", err)
			continue
		}

		switch readString {
		case "1":
			showPrivateChatHelp()

		case "2":
			boo, err = Online(boo)
			if err != nil {
				return err
			}
			continue

		case "3":
			boo, err = UnderLine(boo)
			if err != nil {
				return err
			}
			continue

		case "4":
			err1 := List()
			if err1 != nil {
				fmt.Println("显示所有用户失败：", err1)
				return err
			}

		case "5":
			err := dial.Close()
			if err != nil {
				fmt.Println("关闭与服务端连接失败")
				return err
			}
			fmt.Println("已退出聊天室")
			return nil

		default:
			err2 := WriteTO(boo, readString)
			if err2 != nil {
				return err
			}
		}

	}
}

// Online 上线功能
func Online(boo bool) (bool, error) {
	if boo {
		fmt.Println("已上线，请勿重复上线功能")
	} else {
		fmt.Println("已上线，可以与别人交流")

		marshal, err2 := same.CreateMessage("Online", "")
		if err2 != nil {
			fmt.Println("上线功能序列化失败：", err2)
			return false, err2
		}
		err := same.Write(marshal, dial)
		if err != nil {
			fmt.Println("上线功能写入出现错误", err)
			return false, err
		}

		return true, nil
	}
	return true, nil
}

// UnderLine 下线功能
func UnderLine(boo bool) (bool, error) {
	if !boo {
		fmt.Println("已下线，请勿重复上线功能")
	} else {
		fmt.Println("已下线，无法回复消息")

		marshal, err2 := same.CreateMessage("UnderLine", "")
		if err2 != nil {
			fmt.Println("上线功能序列化失败：", err2)
			return false, err2
		}
		err := same.Write(marshal, dial)

		if err != nil {
			fmt.Println("下线功能写入出现错误", err)
			return false, err
		}
		boo = false
		return boo, nil
	}
	return false, nil
}

// WriteTO 发送功能
func WriteTO(boo bool, readString string) error {
	if !boo {
		fmt.Println("正在下线状态无法发送信息")
		return nil
	}
	var marshal string
	var err1 error

	//判断是私发还是群发功能
	if len(readString) >= 2 && readString[0:2] == "TO" {

		marshal, err1 = same.CreateMessage("Private", readString)
		if err1 != nil {
			fmt.Println("客户端私发消息错误：", err1)
			return err1
		}
	} else {
		marshal, err1 = same.CreateMessage("Every", readString)
		if err1 != nil {
			fmt.Println("客户端群发消息错误：", err1)
			return err1
		}
	}

	err := same.Write(marshal, dial)
	if err != nil {
		fmt.Println("发送消息功能写入出现错误")
		return err
	}
	return nil
}

// List 列表功能
func List() error {

	marshal, err1 := same.CreateMessage("List", "")
	if err1 != nil {
		fmt.Println("客户端群发消息错误：", err1)
		return err1
	}

	err := same.Write(marshal, dial)
	if err != nil {
		fmt.Println("列表功能写入出现错误", err)
		return err
	}
	return nil
}
