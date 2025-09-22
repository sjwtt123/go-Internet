package main

import (
	"fmt"
	method "go-Internet/tcp/method/Clientmethod"
	"net"
)

func main() {

	dial, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}
	fmt.Println("客户端创建成功")

	method.Createprocess(dial)

}
