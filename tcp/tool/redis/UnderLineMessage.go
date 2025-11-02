package redis

import (
	"fmt"
	same "go-Internet/tcp/Samemethod"
)

func ReceiveHistoryAllMessage(username string) error {

	result, err := Rdb.LRange(ctx, "History", 0, -1).Result()
	if err != nil {
		return fmt.Errorf("从redis获取历史消息失败：%v", err)
	}

	if len(result) == 0 {
		fmt.Printf("用户 %s 没有离线消息\n", username)
		return nil
	}

	for _, s := range result {
		from, to, content, err1 := same.AnalyzeHistoryMessage(s)
		if err1 != nil {
			return fmt.Errorf("接受历史消息中：%v", err1)
		}

		if to == "" {
			fmt.Printf("群发消息[%v]：%v\n", from, content)
		} else if to == username || from == username {
			fmt.Printf("私发消息[%v]to[%v]:%v\n", from, to, content)
		}

	}
	return nil
}

func AddHistoryToList(message string) error {
	err := Rdb.RPush(ctx, "History", message).Err()
	if err != nil {
		return fmt.Errorf("存入数据库的历史消息出错：%v", err)
	}
	Rdb.LTrim(ctx, "History", -100, -1)
	return nil
}
