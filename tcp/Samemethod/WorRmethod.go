package Samemethod

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

func Write(msg string, conn net.Conn) error {
	data := []byte(msg)

	// 写长度前缀（4 字节，大端序）
	length := uint32(len(data))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return err
	}
	buf.Write(data)

	// 一次发送整个缓冲区
	_, err = conn.Write(buf.Bytes())
	return err
}
func Read(conn net.Conn) (string, error) {
	// 先读 4 字节长度
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return "", err
	}

	msgLen := binary.BigEndian.Uint32(lenBuf)

	// 再读消息体
	msg := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msg); err != nil {
		return "", err
	}

	return string(msg), nil
}
