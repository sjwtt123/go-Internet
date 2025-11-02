package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	MessageStream = "user_messages" // 统一的消息流
)

// ClientSendMessage 发送消息到目标流中
func ClientSendMessage(from string, content string) error {
	values := map[string]interface{}{
		"from_user_id": from,
		"content":      content,
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
func ReceiveRealtimeMessage(username string) (string, error) {
	lastID := "$" // 从最新消息开始（跳过旧消息）

	// 阻塞读取（一直等待直到有新消息）
	streams, err := Rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{MessageStream, lastID},
		Block:   0, // 0 表示永久阻塞
		Count:   1, // 每次取一条
	}).Result()

	if err != nil {
		return "", fmt.Errorf("读取实时消息失败：%v\n", err)
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			lastID = message.ID // 更新读取位置（下次从此之后读）
			from := message.Values["from_user_id"]
			content := message.Values["content"]
			if from != username {
				return "", nil
			}
			return content.(string), nil
		}
	}
	return "", nil
}
