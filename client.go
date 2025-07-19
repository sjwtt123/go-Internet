package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var isCreate bool

func main() {

	dial, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}
	fmt.Println("客户端创建成功")

	go read(dial)
	write(dial)

	defer func(dial net.Conn) {
		err := dial.Close()
		if err != nil {
			fmt.Println("关闭失败")
		}
	}(dial)
}

func read(dial net.Conn) {
	for {
		by := make([]byte, 1024)
		i, err := dial.Read(by)
		isCreate = true
		if err != nil {
			err := dial.Close()
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("%s", string(by)[:i])
	}
}
func write(dial net.Conn) {
	fmt.Println("------欢迎来到公共网络聊天室------")
	fmt.Println("------1.私聊功能说明------")
	fmt.Println("------2.上线功能------")
	fmt.Println("------3.下线功能------")
	fmt.Println("------4.显示所有用户------")
	fmt.Println("------5.退出登录------")
	boo := true
	for {

		//从终端取到数据
	LOOP:
		reader := bufio.NewReader(os.Stdin)
		readString, err1 := reader.ReadString('\n')
		readString = strings.TrimSpace(readString)

		if !isCreate {
			fmt.Println("未给该客户创建用户名")
			goto LOOP
		}
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
			boo = underLine(dial, boo)
			continue
		case "4":
			err1 := list(dial)
			if err1 != nil {
				fmt.Println("显示所有用户失败：", err1)
			}
		case "5":
			fmt.Println("已退出聊天室")
			return
		default:
			err := writeTO(dial, boo, readString)
			if err != nil {
				fmt.Println("客户端聊天发送数据失败")
			}
		}

	}
}
func Online(dial net.Conn, boo bool) bool {
	if boo {
		fmt.Println("已上线，请勿重复上线功能")
	} else {
		fmt.Println("已上线，可以与别人交流")
		_, err := dial.Write([]byte("re"))
		if err != nil {
			fmt.Println("re写入数据失败:", err)
		}
		return true
	}
	return true
}
func underLine(dial net.Conn, boo bool) bool {
	if !boo {
		fmt.Println("已下线，请勿重复上线功能")
	} else {
		fmt.Println("已下线，无法回复消息")
		_, err := dial.Write([]byte("exit"))
		if err != nil {
			fmt.Println("下线数据写入错误:", err)
		}
		boo = false
		return boo
	}
	return false
}
func writeTO(dial net.Conn, boo bool, readString string) error {
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
func list(dial net.Conn) error {
	_, err2 := dial.Write([]byte("LIST"))
	if err2 != nil {
		return err2
	}
	return nil
}
