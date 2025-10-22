package Samemethod

import (
	"encoding/json"
	"fmt"
)

type SendMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Marshal 将信息序列化
func (receiver *SendMessage) Marshal() (string, error) {
	marshal, err := json.Marshal(receiver)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// Unmarshal 将信息反序列化
func (receiver *SendMessage) Unmarshal(str string) (error, *SendMessage) {
	err := json.Unmarshal([]byte(str), receiver)
	if err != nil {
		return err, nil
	}
	return err, receiver
}

// CreateMessage 通过传入消息类型和消息来创建消息序列
func CreateMessage(ty string, message string) (string, error) {

	sendMessage := &SendMessage{
		Type:    ty,
		Message: message,
	}

	marshal, err := sendMessage.Marshal()
	if err != nil {
		return "", err
	}
	return marshal, nil

}

// AnalyzeMessage 通过将消息反序列化后得到消息类型和数据
func AnalyzeMessage(message string) (string, string, error) {

	sendMessage := &SendMessage{
		Type:    "",
		Message: "",
	}
	err, s := sendMessage.Unmarshal(message)
	if err != nil {
		fmt.Println("接受消息处理反序列化失败")
		return "", "", err
	}
	return s.Message, s.Type, nil
}
