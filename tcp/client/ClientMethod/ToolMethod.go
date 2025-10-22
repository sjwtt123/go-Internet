package ClientMethod

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetfromShell 从终端取到数据
func GetfromShell() (string, error) {
	readString, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("从终端获取数据错误：", err)
		return "", err
	}
	st := strings.TrimSpace(readString)
	return st, err
}

// showWelcomeMessage 显示欢迎信息
func showWelcomeMessage() {
	fmt.Println("------欢迎来到公共网络聊天室------")
	fmt.Println("1. 私聊功能说明")
	fmt.Println("2. 上线功能")
	fmt.Println("3. 下线功能")
	fmt.Println("4. 显示所有用户")
	fmt.Println("5. 退出登录")
	fmt.Println("--------------------------------")
}

// showPrivateChatHelp 显示私聊帮助
func showPrivateChatHelp() {
	fmt.Println("------私聊功能说明------")
	fmt.Println("私发格式：TO[用户名]:内容")
	fmt.Println("示例：TO[alice]:你好")
	fmt.Println("----------------------")
}
