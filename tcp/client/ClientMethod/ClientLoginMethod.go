package ClientMethod

import (
	"fmt"
	same "go-Internet/tcp/ReadWritermethod"
	"net"
)

// Createprocess 创建登录或注册进程
func Createprocess(conn net.Conn) error {
	dial = conn
	for {
		fmt.Println("请选择登录：1or注册：2")
		selectString, err := GetfromShell()
		if err != nil {
			fmt.Println("登录，注册功能从终端获取数据失败", err)
			return err
		}

		switch selectString {
		case "1":
			err := LoginUser()
			if err != nil {
				return err
			}
			return nil
		case "2":
			err := RegisterUser()
			if err != nil {
				return err
			}
			return nil
		default:
			fmt.Println("请输入正确的选择")
			continue
		}
	}
}

// RegisterUser 注册新用户
func RegisterUser() error {
	err := same.Write("2", dial)
	if err != nil {
		return err
	}

	for {
		username, passwd := GetNameAndPasswd()

		err := same.Write(username, dial)
		if err != nil {
			fmt.Println("注册时写入用户名数据错误", err)
			return err
		}
		scanner, err := same.Read(dial)
		if err != nil {
			fmt.Println("读入用户名是否重复数据失败", err)
			return err
		}

		for scanner.Scan() {
			if scanner.Text() == "isCreate" {
				fmt.Println("用户名已存在，请重新输入")
				break
			} else {
				err1 := same.Write(passwd, dial)
				if err1 != nil {
					fmt.Println("传入密码数据失败", err1)
					return err1
				}
				return nil
			}
		}

	}

}

// LoginUser 用户登录
func LoginUser() error {
	err := same.Write("1", dial)
	if err != nil {
		return err
	}

LOOP:
	for {
		username, passwd := GetNameAndPasswd()
		err := same.Write(username, dial)
		if err != nil {
			fmt.Println("登录时传入数据错误", err)
			return err
		}
		scanner, err := same.Read(dial)
		if err != nil {
			fmt.Println("登录时从服务端读入用户名验证失败", err)
			return err
		}

		for scanner.Scan() {
			if scanner.Text() == "noCreate" {
				fmt.Println("用户名不存在或已登录，请重新输入")
				goto LOOP
			} else {
				err := same.Write(passwd, dial)
				if err != nil {
					return fmt.Errorf("%s传入密码数据失败", err)
				}
				break
			}
		}

		scanner1, err := same.Read(dial)
		if err != nil {
			fmt.Println("登录时从服务端读入密码验证失败", err)
			return err
		}

		for scanner1.Scan() {
			if scanner1.Text() == "false" {
				fmt.Println("密码错误，请重新输入")
				break
			} else {
				fmt.Println("登录成功")
				return nil
			}
		}

	}
}

// GetNameAndPasswd 从终端获取用户名和密码
func GetNameAndPasswd() (string, string) {
	for {
		fmt.Println("请输入用户名：")
		username, err1 := GetfromShell()

		fmt.Println("请输入密码：")
		passwd, err2 := GetfromShell()

		if err2 != nil || err1 != nil {
			fmt.Println("写入数据不正确请重新输入", err2)
			continue
		}
		if username == "" || passwd == "" {
			fmt.Println("用户名或密码不能为空")
			continue
		}
		return username, passwd
	}
}
