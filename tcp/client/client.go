package main

import (
	"context"
	"fmt"
	method "go-Internet/tcp/client/ClientMethod"
	"net"
	"runtime/debug"
)

func main() {
	// 设置主协程的 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("主协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	//进行消息判断业务处理
	go func() {
		err := method.Start()
		if err != nil {
			cancel()
		}
	}()

	//开启读协程一直并发在读
	go func() {
		err := method.Read(dial)
		if err != nil {
			cancel()
		}
	}()

	//开启心跳检测判断客户端还在连接
	go func() {
		err := method.WriteHeart()
		if err != nil {
			cancel()
		}
	}()

	<-ctx.Done()
	fmt.Println("程序退出")
}
