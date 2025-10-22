package ServeMethod

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
	"go-Internet/tcp/tool/mysql"
	"net"
	"runtime/debug"
)

// Init 初始化mysql
func Init() {
	//初始化数据库
	err, dbs := mysql.Start()
	db = dbs
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}

}

// ISLoginOrCreate 判断登录或者创建用户
func ISLoginOrCreate(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ISLoginOrCreate()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	Init()

	scanner, err := same.Read(conn)
	if err != nil {
		fmt.Println("判断登录或者注册功能读取失败：", err)
		return
	}

	for scanner.Scan() {
		ion := scanner.Text()
		if ion == "2" {
			err := ReIsHaveUser(conn)
			if err != nil {
				fmt.Println("注册出现错误", err)
				return
			}
		} else {
			err := LoIsNotHaveUser(conn)
			if err != nil {
				fmt.Println("登录出现错误", err)
				return
			}
		}

	}

}

// LoIsNotHaveUser 登录功能
func LoIsNotHaveUser(conn net.Conn) error {
	var passwd, username string
LOOP:
	//获取用户名
	scanner, err := same.Read(conn)
	if err != nil {
		return fmt.Errorf("登录时读出数据失败:%v", err)
	}

	for scanner.Scan() {
		username = scanner.Text()
		//查询数据库是否存在
		isHaveUser := IsHaveUser(username)

		isHaveClient := FindClient(username)

		if !isHaveUser || isHaveClient.Nickname != "" {
			err2 := same.Write("noCreate", conn)
			if err2 != nil {
				return fmt.Errorf("登录判断用户存在写入数据失败:%v", err2)
			}
			goto LOOP
		} else {
			//用户校验成功返回
			err2 := same.Write("success", conn)
			if err2 != nil {
				return fmt.Errorf("登录成功写入数据失败:%v", err2)
			}

			//获取密码
			scanner1, err1 := same.Read(conn)
			if err1 != nil {
				return fmt.Errorf("登录获取密码数据失败:%v", err1)
			}

			for scanner1.Scan() {
				passwd = scanner1.Text()
			}
			//判断密码是否正确
			boo, err := IsCurrentPasswd(username, passwd, conn)
			if err != nil {
				return fmt.Errorf("登录密码校验失败:%v", err)
			}
			if !boo {
				goto LOOP
			}
		}
	}
	Createprocess(username, conn)

	return nil
}

// IsCurrentPasswd 校验密码
func IsCurrentPasswd(username string, passwd string, conn net.Conn) (bool, error) {
	namePasswd, err := mysql.SelectByNamePasswd(db, username, passwd)
	if err != nil {
		return false, err
	}
	if namePasswd {
		err = same.Write("true", conn)
		return true, err

	} else {
		err = same.Write("false", conn)
		return false, err

	}

}

// ReIsHaveUser 注册功能
func ReIsHaveUser(conn net.Conn) error {
	var passwd string

	//获取用户名
LOOP:
	scanner, err := same.Read(conn)
	if err != nil {
		return fmt.Errorf("注册读入用户名出错：%v", err)
	}

	for scanner.Scan() {
		username := scanner.Text()
		name := IsHaveUser(username)

		if name {
			err1 := same.Write("isCreate", conn)
			if err1 != nil {
				return fmt.Errorf("注册判断重复用户写入数据失败:%v", err1)
			}
			goto LOOP

		} else {
			err1 := same.Write("success", conn)
			if err1 != nil {
				return fmt.Errorf("注册判断重复用户写入数据失败:%v", err1)
			}

			//获取密码
			scanner1, err2 := same.Read(conn)
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

			err4 := same.Write(username+"用户注册成功", conn)
			if err4 != nil {
				return fmt.Errorf("用户注册成功信息写入失败:%v", err4)
			}

			Createprocess(username, conn)

		}
	}

	return nil
}

// Createprocess 创建用户实例
func Createprocess(username string, conn net.Conn) {

	//创建用户对象存每个客户的用户名，加入到在线客户端中
	client := Client{Nickname: username, Conn: conn, Boo: true}
	HbManager.AddClient(&client)

	messageChan <- fmt.Sprintf("%s进入聊天室", username)

	//开启信息处理流程
	CRead(conn, client)

}
