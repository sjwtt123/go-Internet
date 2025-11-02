package Samemethod

import (
	"encoding/json"
	"fmt"
)

type SendMessage struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Content    string `json:"content"`
	Tp         string `json:"tp"`
}

// Marshal 将信息序列化
func (receiver *SendMessage) Marshal() (string, error) {
	marshal, err := json.Marshal(receiver)
	if err != nil {
		return "", fmt.Errorf("将发送消息序列化错误：%v", err)
	}
	return string(marshal), nil
}

// Unmarshal 将信息反序列化
func (receiver *SendMessage) Unmarshal(str string) (error, *SendMessage) {
	err := json.Unmarshal([]byte(str), receiver)
	if err != nil {
		return fmt.Errorf("将发送消息反序列化错误：%v", err), nil
	}
	return nil, receiver
}

// CreateMessage 创建消息序列
func CreateMessage(from string, to string, ty string, message string) (string, error) {

	sendMessage := &SendMessage{
		FromUserId: from,
		ToUserId:   to,
		Tp:         ty,
		Content:    message,
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
		FromUserId: "",
		ToUserId:   "",
		Tp:         "",
		Content:    message,
	}
	err, s := sendMessage.Unmarshal(message)
	if err != nil {
		return "", "", err
	}
	return s.Content, s.Tp, nil
}

// AnalyzeHistoryMessage 反序列化得到需要的历史数据
func AnalyzeHistoryMessage(message string) (string, string, string, error) {

	sendMessage := &SendMessage{
		FromUserId: "",
		ToUserId:   "",
		Tp:         "",
		Content:    message,
	}
	err, s := sendMessage.Unmarshal(message)
	if err != nil {
		return "", "", "", err
	}

	return s.FromUserId, s.ToUserId, s.Content, err

}
