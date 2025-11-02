package main

import (
	"context"
	"fmt"
	method "go-Internet/tcp/serve/ServeMethod"
	"log"
	"net"
	"runtime/debug"
)

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		fmt.Println("创建服务端失败", err)
	}
	fmt.Println("创建服务端成功，等待连接")

	// 设置主协程的 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("主协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	//开启心跳检测
	go method.HbManager.Start()

	for {

		//等待下一个连接到该接口的连接
		accept, err2 := listen.Accept()
		if err2 != nil {
			log.Println("客户端连接失败：", err2)
			break
		}
		fmt.Printf("连接成功对应的ip地址，端口号：%v\n", accept.RemoteAddr().String())

		// 为每个客户端创建一个独立的 context 控制生命周期
		ctx, cancel := context.WithCancel(context.Background())

		client := &method.Client{
			Conn:     accept,
			Nickname: "",
			Boo:      true,
		}

		//为每个用户开启广播协程
		go method.Radio(ctx)

		//开启登录功能
		go func() {
			err = client.ISLoginOrCreate()
			if err != nil {
				cancel()
			}
			//开启信息处理流程
			go func() {
				err = client.CRead()
				if err != nil {
					cancel()
				}
			}()

		}()

	}
}
