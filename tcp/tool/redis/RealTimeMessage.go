package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	MessageStream = "user_messages" // 统一的消息流
)

var MessageHandler func(msg Message)

type Message struct {
	Content string
	Type    string
	From    string
}

// ClientSendMessage 发送消息到目标流中
func ClientSendMessage(content string, ty string, form string) error {
	values := map[string]interface{}{
		"content": content,
		"type":    ty,
		"form":    form,
	}
	err := Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: MessageStream,
		MaxLen: 100,
		Limit:  0,
		Values: values,
		Approx: true,
	}).Err()
	if err != nil {

		return fmt.Errorf("写入数据流失败：%v", err)
	}
	return nil
}

// ReceiveRealtimeMessage 使用 XRead 实时监听消息（非消费者组版）
func ReceiveRealtimeMessage() {
	lastID := "$"

	for {
		streams, err := Rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{MessageStream, lastID},
			Block:   0,
			Count:   1,
		}).Result()

		if err != nil {
			fmt.Printf("读取实时消息失败：%v\n", err)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				lastID = message.ID

				msg := Message{
					Content: message.Values["content"].(string),
					Type:    message.Values["type"].(string),
					From:    message.Values["form"].(string),
				}

				if MessageHandler != nil {
					MessageHandler(msg)
				}
			}
		}
	}
}
