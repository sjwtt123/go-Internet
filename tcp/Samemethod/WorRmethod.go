package Samemethod

import (
	"bufio"
	"net"
	"strings"
)

var by = make([]byte, 1024)

func Read(conn net.Conn) (*bufio.Scanner, error) {
	i, err := conn.Read(by)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(by[:i])))
	return scanner, nil
}
func Write(s string, conn net.Conn) error {
	s += "\n"
	_, err := conn.Write([]byte(s))
	if err != nil {
		return err
	}
	return err
}
