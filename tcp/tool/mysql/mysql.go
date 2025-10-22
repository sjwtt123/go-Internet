package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Passwd string `db:"passwd"`
}

func Start() (error, *sqlx.DB) {

	open, err := sqlx.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/system1")
	if err != nil {
		return fmt.Errorf("%s打开MySQL失败", err), open
	}

	return nil, open
}

func SelectOneByName(db *sqlx.DB, name string) (bool, error) {
	var users []User
	err := db.Select(&users, "select * from tcp_user where name=?", name)
	if err != nil {
		return false, nil
	}
	if len(users) == 1 {
		return true, err
	}
	return false, nil
}

func SelectByNamePasswd(db *sqlx.DB, name string, passwd string) (bool, error) {

	var users []User
	err := db.Select(&users, "select * from tcp_user where name=? and passwd=?", name, passwd)
	if err != nil {
		return false, nil
	}
	if len(users) == 1 {
		return true, err
	}
	return false, nil

}

func InsertDate(u *User, db *sqlx.DB) error {
	_, err := db.Exec("insert into tcp_user(name,passwd) values (?,?)", u.Name, u.Passwd)
	if err != nil {
		return err
	}
	return nil
}
