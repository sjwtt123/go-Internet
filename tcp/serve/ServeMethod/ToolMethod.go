package ServeMethod

import (
	"fmt"
	"go-Internet/tcp/tool/mysql"
	"strings"
)

// IsHaveUser 查询数据库中是否有该用户
func IsHaveUser(username string) bool {

	boo, err := mysql.SelectOneByName(db, username)
	if err != nil {
		fmt.Println("查询出错：", err)
	}
	return boo

}

// FindName 私发功能找用户名
func FindName(s string) (string, string, error) {
	index := strings.Index(s, "[")
	index1 := strings.Index(s, "]")
	if index < 0 || index1 < 0 || index1+2 >= len(s) {
		return "", "", fmt.Errorf("输入格式错误\n")
	}
	message := s[index1+2:]
	s = s[index+1 : index1]
	return s, message, nil
}

// FindClient 寻找用户
func FindClient(name string) Client {
	clientManager.Mutex.Lock()
	defer clientManager.Mutex.Unlock()
	for client, _ := range clientManager.clients {
		if client.nickname == name {
			return client
		}
	}

	return Client{}

}

// CloseConn 关闭与客户端连接
func CloseConn(client Client) {
	err := client.conn.Close()
	if err != nil {
		fmt.Printf("关闭%v客户端失败%v", client.nickname, err)
		return
	}

}
