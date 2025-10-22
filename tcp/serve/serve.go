package main

import (
	"fmt"
	method "go-Internet/tcp/serve/ServeMethod"
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

	//延迟关闭
	defer func(listen net.Listener) {
		err1 := listen.Close()
		if err1 != nil {
			fmt.Println("listen关闭失败", err1)
		}
		//返回该接口的网络地址
		fmt.Println(listen.Addr().String())
	}(listen)

	//开启心跳检测
	go method.HbManager.Start()

	for {

		//等待下一个连接到该接口的连接
		accept, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println("客户端连接失败：", err2)
			break
		}

		fmt.Printf("连接成功对应的ip地址，端口号：%v\n", accept.RemoteAddr().String())

		//为每个用户开启广播协程
		go method.Radio()

		//开启登录功能
		go method.ISLoginOrCreate(accept)

	}
}
