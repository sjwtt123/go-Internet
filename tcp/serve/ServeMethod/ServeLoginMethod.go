package ServeMethod

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
	"go-Internet/tcp/tool/mysql"
	"go-Internet/tcp/tool/redis"
	_ "go-Internet/tcp/tool/redis"
	"log"
	"net"
	"runtime/debug"
)

// init 初始化mysql
func init() {
	//初始化数据库
	err, dbs := mysql.Start()
	db = dbs
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}

}

// ISLoginOrCreate 判断登录或者创建用户
func (client *Client) ISLoginOrCreate() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ISLoginOrCreate()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()
	for {
		scanner, err := same.Read(client.Conn)
		if err != nil {
			return fmt.Errorf("判断登录或者注册功能读取失败：%v", err)
		}

		for scanner.Scan() {
			ion := scanner.Text()
			if ion == ServeCommandRegister {
				err = client.ReIsHaveUser()
				if err != nil {
					if err.Error() == "注册失败" {
						break
					}
					log.Println(err)
					return err

				}
			} else if ion == ServeCommandLogin {
				err = client.LoIsNotHaveUser()
				if err != nil {
					if err.Error() == "登录失败" {
						break
					}
					log.Println(err)
					return err

				}
			}
			return nil
		}

	}
}

// LoIsNotHaveUser 登录功能
func (client *Client) LoIsNotHaveUser() error {
	var passwd, username string

	//获取用户名
	scanner, err := same.Read(client.Conn)
	if err != nil {
		return fmt.Errorf("登录时读出数据失败:%v", err)
	}

	for scanner.Scan() {
		username = scanner.Text()
		//查询数据库是否存在
		isHaveUser := IsHaveUser(username)

		isHaveClient := FindClient(username)

		if !isHaveUser || isHaveClient.Nickname != "" {
			err2 := same.Write(ReceiveUserExistsOrLo, client.Conn)
			if err2 != nil {
				return fmt.Errorf("登录判断用户存在写入数据失败:%v", err2)
			}

			return fmt.Errorf("登录失败")

		} else {
			//用户校验成功返回
			err2 := same.Write(ReceiveSuccess, client.Conn)
			if err2 != nil {
				return fmt.Errorf("登录成功写入数据失败:%v", err2)
			}

			//获取密码
			scanner1, err1 := same.Read(client.Conn)
			if err1 != nil {
				return fmt.Errorf("登录获取密码数据失败:%v", err1)
			}

			for scanner1.Scan() {
				passwd = scanner1.Text()
			}
			//判断密码是否正确
			boo, err := IsCurrentPasswd(username, passwd, client.Conn)
			if err != nil {
				return fmt.Errorf("登录密码校验失败:%v", err)
			}
			if !boo {
				return fmt.Errorf("登录失败")
			}
			client.Nickname = username
			client.CreateprocessAndStart()
		}
	}

	return nil
}

// IsCurrentPasswd 校验密码
func IsCurrentPasswd(username string, passwd string, conn net.Conn) (bool, error) {
	namePasswd, err := mysql.SelectByNamePasswd(db, username, passwd)
	if err != nil {
		return false, err
	}
	if namePasswd {
		err = same.Write(ReceiveTrue, conn)
		return true, err

	} else {
		err = same.Write(ReceiveFalse, conn)
		return false, err

	}

}

// ReIsHaveUser 注册功能
func (client *Client) ReIsHaveUser() error {
	var passwd string

	//获取用户名
	scanner, err := same.Read(client.Conn)
	if err != nil {
		return fmt.Errorf("注册读入用户名出错：%v", err)
	}

	for scanner.Scan() {
		username := scanner.Text()

		//判断是否存在用户
		isExistname := IsHaveUser(username)

		if isExistname {
			err1 := same.Write(ReceiveUserExists, client.Conn)
			if err1 != nil {
				return fmt.Errorf("注册判断重复用户写入数据失败:%v", err1)
			}
			fmt.Println("用户已存在请重新创建")
			return fmt.Errorf("注册失败")
		} else {

			err1 := same.Write(ReceiveSuccess, client.Conn)
			if err1 != nil {
				return fmt.Errorf("注册判断重复用户写入数据失败:%v", err1)
			}

			//获取密码
			scanner1, err2 := same.Read(client.Conn)
			if err2 != nil {
				return fmt.Errorf("注册获取密码读取数据失败:%v", err2)
			}
			for scanner1.Scan() {
				passwd = scanner1.Text()
			}

			//将用户名，密码存入数据库中
			err3 := mysql.InsertDate(&mysql.User{Name: username, Passwd: passwd}, db)
			if err3 != nil {
				return fmt.Errorf("注册新用户存入数据库失败：%v", err3)
			}

			err4 := same.Write(username+"用户注册成功", client.Conn)
			if err4 != nil {
				return fmt.Errorf("用户注册成功信息写入失败:%v", err4)
			}

			//注册后将用户写入redis，初始活跃度为1
			redis.AddIntoActive(username, 1)

			client.Nickname = username
			client.CreateprocessAndStart()
		}
	}

	return nil
}

// CreateprocessAndStart  创建用户实例
func (client *Client) CreateprocessAndStart() {

	//创建用户对象存每个客户的用户名，加入到在线客户端中
	HbManager.AddClient(client)
	messageChan <- fmt.Sprintf("%s进入聊天室", client.Nickname)

}
