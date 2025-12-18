package main

import (
	"context"
	"fmt"
	method "go-Internet/tcp/client/ClientMethod"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func main() {
	// 设置主协程的 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			log.Printf("主协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dial, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Println("连接失败：", err)
		return
	}

	//创建用户实例
	client := &method.Client{
		Boo:      true,
		Dial:     dial,
		Nickname: "",
	}
	fmt.Println("客户端创建成功")

	//进行用户信息登录或者注册
	err1 := client.Createprocess()
	if err1 != nil {
		log.Println(err1)
		return
	}

	//开启心跳检测判断客户端还在连接
	go func() {
		err = client.WriteHeart()
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	//进行消息判断业务处理
	go func() {
		err = client.Start()
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	//开启读协程一直并发在读
	go func() {
		err = client.Read()
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	//  捕获系统退出信号（Ctrl+C、kill、终端关闭等）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("逻辑触发退出")
	case sig := <-sigChan:
		log.Printf("收到系统信号: %v\n", sig)
		client.Leave()
	}

	fmt.Println("程序退出")
}
