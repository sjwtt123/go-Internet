package main

import (
	"fmt"
	method "go-Internet/tcp/client/ClientMethod"
	"net"
)

func main() {
	dial, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}
	fmt.Println("客户端创建成功")

	//进行用户信息登录或者注册
	err1 := method.Createprocess(dial)
	if err1 != nil {
		return
	}

	//开启读协程一直并发在读
	go method.Read(dial)

	//进行消息判断业务处理
	method.Start()

}
