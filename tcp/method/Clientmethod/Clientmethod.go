package Clientmethod

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	boo = true
	by  = make([]byte, 1024)
)

func Read(dial net.Conn) {
	scanner := bufio.NewScanner(dial)
	for scanner.Scan() {
		msg := scanner.Text() // 自动按\n拆分
		fmt.Println(msg)
	}

	// 检查退出原因：错误或EOF
	if err := scanner.Err(); err != nil {

		fmt.Println("连接关闭，读取结束")
	}
}

// Createprocess 创建新用户信息
func Createprocess(conn net.Conn) {

LOOP:
	fmt.Println("请创建用户名：")
	readName, err := bufio.NewReader(os.Stdin).ReadString('\n')
	username := strings.TrimSpace(readName)
	if username == "" {
		fmt.Println("用户名不能为空")
		goto LOOP
	}

	conn.Write([]byte(username))
	i, err := conn.Read(by)
	if err != nil {
		fmt.Println("读入失败，客户端已关闭")
	}

	scanner := bufio.NewScanner(strings.NewReader(string(by[:i])))
	for scanner.Scan() {
		if scanner.Text() == "isCreate" {
			fmt.Println("用户名已存在")
			goto LOOP
		}

		if err != nil {
			fmt.Println("输入用户名错误：", err)
		}
		fmt.Println("创建成功")

	}

	Write(conn)

}
func Write(dial net.Conn) {
	go Read(dial)
	fmt.Println("------欢迎来到公共网络聊天室------")
	fmt.Println("------1.私聊功能说明------")
	fmt.Println("------2.上线功能------")
	fmt.Println("------3.下线功能------")
	fmt.Println("------4.显示所有用户------")
	fmt.Println("------5.退出登录------")

	for {
		//从终端取到数据
		reader := bufio.NewReader(os.Stdin)
		readString, err1 := reader.ReadString('\n')
		readString = strings.TrimSpace(readString)

		if err1 != nil {
			fmt.Println("写入聊天室数据失败")
		}
		switch readString {
		case "1":
			fmt.Println("私发方式：TO [用户名]:内容")
		case "2":
			boo = Online(dial, boo)
			continue
		case "3":
			boo = UnderLine(dial, boo)
			continue
		case "4":
			err1 := List(dial)
			if err1 != nil {
				fmt.Println("显示所有用户失败：", err1)
			}
		case "5":
			fmt.Println("已退出聊天室")
			return
		default:
			err := WriteTO(dial, boo, readString)
			if err != nil {
				fmt.Println("客户端聊天发送数据失败")
			}
		}

	}
}

// Online 上线功能
func Online(dial net.Conn, boo bool) bool {
	if boo {
		fmt.Println("已上线，请勿重复上线功能")
	} else {
		fmt.Println("已上线，可以与别人交流")
		_, err := dial.Write([]byte("RE\n"))
		if err != nil {
			fmt.Println("re写入数据失败:", err)
		}
		return true
	}
	return true
}

// UnderLine 下线功能
func UnderLine(dial net.Conn, boo bool) bool {
	if !boo {
		fmt.Println("已下线，请勿重复上线功能")
	} else {
		fmt.Println("已下线，无法回复消息")
		_, err := dial.Write([]byte("EXIT\n"))
		if err != nil {
			fmt.Println("下线数据写入错误:", err)
		}
		boo = false
		return boo
	}
	return false
}

// WriteTO 私发功能
func WriteTO(dial net.Conn, boo bool, readString string) error {
	if !boo {
		fmt.Println("正在下线状态无法发送信息")
		return nil
	}
	_, err2 := dial.Write([]byte(readString))
	if err2 != nil {
		return err2
	}
	return nil
}

// List 列表功能
func List(dial net.Conn) error {
	_, err2 := dial.Write([]byte("LIST\n"))
	if err2 != nil {
		return err2
	}
	return nil
}
