package ServeMethod

import (
	"context"
	"fmt"
	"go-Internet/tcp/tool/mysql"
	"go-Internet/tcp/tool/redis"
	"log"
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

func (client *Client) HandleLogic(ctx context.Context) error {
	if client.LoginOrCreate() != nil {
		return fmt.Errorf("登录注册功能失败")
	}

	for {

		select {
		case <-ctx.Done():
			return nil
		case msg := <-client.Inchan:
			err := client.SRead(msg)
			if err != nil {
				return err
			}
		}

	}

}

// LoginOrCreate 判断登录或者创建用户
func (client *Client) LoginOrCreate() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ISLoginOrCreate()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()
	for {
		msg := <-client.Inchan

		if msg == ServeCommandRegister {
			err := client.RegisterUser()
			if err != nil {
				if err.Error() == "注册失败" {
					continue
				}
				log.Println(err)
				return err

			}
		} else if msg == ServeCommandLogin {
			err := client.LoginUser()
			if err != nil {
				if err.Error() == "登录失败" {
					continue
				}
				log.Println(err)
				return err

			}
		}
		break
	}

	return nil
}

// LoginUser 登录功能
func (client *Client) LoginUser() error {

	//获取用户名
	username := <-client.Inchan

	//查询数据库是否存在
	isHaveUser := IsHaveUser(username)

	isHaveClient := FindClient(username)

	if !isHaveUser || isHaveClient != nil {
		client.Outchan <- ReceiveUserExistsOrLo

	} else {
		//用户校验成功返回
		client.Outchan <- ReceiveSuccess

		//获取密码
		passwd := <-client.Inchan

		//判断密码是否正确
		boo, err := client.IsCurrentPasswd(username, passwd)
		if err != nil {
			return fmt.Errorf("登录密码校验失败:%v", err)
		}
		if !boo {
			return fmt.Errorf("登录失败")
		}
		client.Nickname = username
		client.CreateprocessAndStart()
	}
	return nil
}

// IsCurrentPasswd 校验密码
func (client *Client) IsCurrentPasswd(username string, passwd string) (bool, error) {
	namePasswd, err := mysql.SelectByNamePasswd(db, username, passwd)
	if err != nil {
		return false, err
	}
	if namePasswd {
		client.Outchan <- ReceiveTrue
		return true, err

	} else {
		client.Outchan <- ReceiveFalse
		return false, err
	}

}

// RegisterUser 注册功能
func (client *Client) RegisterUser() error {

	//获取用户名
	username := <-client.Inchan

	//判断是否存在用户
	isExistname := IsHaveUser(username)

	if isExistname {

		client.Outchan <- ReceiveUserExists
		fmt.Println("用户已存在请重新创建")
		return fmt.Errorf("注册失败")
	} else {
		client.Outchan <- ReceiveSuccess
		//获取密码
		passwd := <-client.Inchan

		//将用户名，密码存入数据库中
		err3 := mysql.InsertDate(&mysql.User{Name: username, Passwd: passwd}, db)
		if err3 != nil {
			return fmt.Errorf("注册新用户存入数据库失败：%v", err3)
		}

		client.Outchan <- username + "用户注册成功"

		//注册后将用户写入redis，初始活跃度为1
		redis.AddIntoActive(username, 1)

		client.Nickname = username
		client.CreateprocessAndStart()
	}
	return nil
}

// CreateprocessAndStart  创建用户实例
func (client *Client) CreateprocessAndStart() {

	//创建用户对象存每个客户的用户名，加入到在线客户端中
	HbManager.AddClient(client)
	MessageChan <- fmt.Sprintf("%s进入聊天室", client.Nickname)

}
