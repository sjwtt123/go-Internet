package main

import (
	"fmt"
	method "go-Internet/tcp/method/Servemethod"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		fmt.Println("创建服务端失败", err)
	}
	fmt.Println("创建服务端成功，等待连接")

	//延迟关闭
	defer func(listen net.Listener) {
		err1 := listen.Close()
		if err1 != nil {
			fmt.Println("listen关闭失败", err1)
		}
		//返回该接口的网络地址
		fmt.Println(listen.Addr().String())
	}(listen)

	for {
		//等待下一个连接到该接口的连接
		accept, err2 := listen.Accept()
		if err2 != nil {

			fmt.Println("客户端连接失败：", err2)
			break
		}
		fmt.Printf("连接成功对应的ip地址，端口号：%v\n", accept.RemoteAddr().String())

		go method.Radio()

		go func() {
			err3 := method.Createprocess(accept)
			if err3 != nil {
				fmt.Println(err3)
			}
		}()

	}
}
