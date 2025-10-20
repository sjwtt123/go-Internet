package ClientMethod

import (
	"bufio"
	"fmt"
	same "go-Internet/tcp/ReadWritermethod"
	"net"
	"os"
	"strings"
)

var (
	boo  = true
	dial net.Conn
)

// Read 并发读协程
func Read(dial net.Conn) {
	for {
		scanner, err := same.Read(dial)
		if err != nil {
			fmt.Println("读取结束,连接关闭")
			break
		}

		// 检查退出原因：错误或EOF
		if err = scanner.Err(); err != nil {
			fmt.Println("连接关闭，读取结束")
			return
		}
		for scanner.Scan() {
			msg := scanner.Text() // 自动按\n拆分
			fmt.Println(msg)
		}

	}
}

// Start 开启业务
func Start() {
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
				return
			}
			continue

		case "3":
			boo, err = UnderLine(boo)
			if err != nil {
				return
			}
			continue

		case "4":
			err1 := List()
			if err1 != nil {
				fmt.Println("显示所有用户失败：", err1)
				return
			}

		case "5":
			err := dial.Close()
			if err != nil {
				fmt.Println("关闭与服务端连接失败")
				return
			}
			fmt.Println("已退出聊天室")
			return

		default:
			err2 := WriteTO(boo, readString)
			if err2 != nil {
				return
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
		err := same.Write("RE", dial)
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
		err := same.Write("EXIT", dial)
		if err != nil {
			fmt.Println("下线功能写入出现错误", err)
			return false, err
		}
		boo = false
		return boo, nil
	}
	return false, nil
}

// WriteTO 私发功能
func WriteTO(boo bool, readString string) error {
	if !boo {
		fmt.Println("正在下线状态无法发送信息")
		return nil
	}
	err := same.Write(readString, dial)
	if err != nil {
		fmt.Println("私发功能写入出现错误")
		return err
	}
	return nil
}

// List 列表功能
func List() error {
	err := same.Write("LIST", dial)
	if err != nil {
		fmt.Println("列表功能写入出现错误", err)
		return err
	}
	return nil
}

// GetfromShell 从终端取到数据
func GetfromShell() (string, error) {
	readString, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("从终端获取数据错误：", err)
		return "", err
	}
	st := strings.TrimSpace(readString)
	return st, err
}

// showWelcomeMessage 显示欢迎信息
func showWelcomeMessage() {
	fmt.Println("------欢迎来到公共网络聊天室------")
	fmt.Println("1. 私聊功能说明")
	fmt.Println("2. 上线功能")
	fmt.Println("3. 下线功能")
	fmt.Println("4. 显示所有用户")
	fmt.Println("5. 退出登录")
	fmt.Println("--------------------------------")
}

// showPrivateChatHelp 显示私聊帮助
func showPrivateChatHelp() {
	fmt.Println("------私聊功能说明------")
	fmt.Println("私发格式：TO[用户名]:内容")
	fmt.Println("示例：TO[alice]:你好")
	fmt.Println("----------------------")
}
