package ClientMethod

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
)

// Createprocess 创建登录或注册进程
func (client *Client) Createprocess() error {
	for {
		fmt.Println("请选择登录：1or注册：2")
		selectString, err := GetfromShell()
		if err != nil {
			return fmt.Errorf("登录，注册功能：%v", err)
		}

		switch selectString {
		case ClientCommandLogin:
			err = client.LoginUser()
			if err != nil {
				if err.Error() == "登录失败" {
					continue
				}
				return err
			}

		case ClientCommandRegister:
			err = client.RegisterUser()
			if err != nil {
				if err.Error() == "注册失败" {
					continue
				}
				return err
			}

		default:
			fmt.Println("请输入正确的选择")
			continue
		}
		return nil
	}
}

// RegisterUser 注册新用户
func (client *Client) RegisterUser() error {
	err := same.Write(ClientCommandRegister, client.Dial)
	if err != nil {
		return fmt.Errorf("传入注册信息失败：%v", err)
	}

	username, passwd := GetNameAndPasswd()

	err = same.Write(username, client.Dial)
	if err != nil {
		return fmt.Errorf("注册时写入用户名数据错误:%v", err)
	}
	scanner, err := same.Read(client.Dial)
	if err != nil {
		return fmt.Errorf("读入用户名是否重复数据失败:%v", err)
	}

	for scanner.Scan() {
		if scanner.Text() == ResponseUserExists {

			fmt.Println("用户已存在请重新输入")
			return fmt.Errorf("注册失败")
		} else {
			err1 := same.Write(passwd, client.Dial)
			if err1 != nil {
				return fmt.Errorf("传入密码数据失败:%v", err1)
			}

			//注册成功赋值nickname
			client.Nickname = username
			return nil
		}
	}

	return nil

}

// LoginUser 用户登录
func (client *Client) LoginUser() error {
	err := same.Write(ClientCommandLogin, client.Dial)
	if err != nil {
		return fmt.Errorf("传入登录信息失败：%v", err)
	}

	username, passwd := GetNameAndPasswd()
	err = same.Write(username, client.Dial)
	if err != nil {
		return fmt.Errorf("登录时传入数据错误：%v", err)
	}

	//验证用户名
	scanner, err := same.Read(client.Dial)
	if err != nil {
		return fmt.Errorf("登录时从服务端读入用户名验证失败:%v", err)
	}

	for scanner.Scan() {
		if scanner.Text() == ResponseUserExistsOrLo {

			fmt.Println("用户不存在或已登录请重新输入")
			return fmt.Errorf("登录失败")
		} else {
			err = same.Write(passwd, client.Dial)
			if err != nil {
				return fmt.Errorf("传入密码数据失败:%v", err)
			}
			break
		}
	}

	//验证密码
	scanner1, err := same.Read(client.Dial)
	if err != nil {
		return fmt.Errorf("登录时从服务端读入密码验证失败:%v", err)
	}

	for scanner1.Scan() {
		if scanner1.Text() == ResponseLoginFailed {
			fmt.Println("登录密码错误,请重新输入")
			return fmt.Errorf("登录失败")
		} else {

			client.Nickname = username
			fmt.Println("登录成功")
			return nil
		}
	}

	return nil
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
