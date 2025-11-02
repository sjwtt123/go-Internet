package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var Rdb *redis.Client
var ctx = context.Background()

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
}

func AddIntoActive(name string, num float64) {
	err := Rdb.ZAdd(ctx, "active", &redis.Z{
		Score:  num,
		Member: name,
	}).Err()
	if err != nil {
		log.Println("redis插入数据失败", err)
		return
	}

}

func IncreaseActive(name string) error {

	err := Rdb.ZIncrBy(ctx, "active", 5, name).Err()
	if err != nil {
		return fmt.Errorf("redis活跃度修改数据失败:%v", err)
	}
	return nil
}

func FindAllActive() ([]redis.Z, error) {

	result, err := Rdb.ZRevRangeWithScores(ctx, "active", 0, -1).Result()

	if err != nil {
		return nil, fmt.Errorf("查询活跃度列表失败：%v", err)
	}

	return result, nil

}
