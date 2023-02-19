package example

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
)

const Ex05UserCountKey = "ex05_user_count"

// Ex05 hash数据结果的运用（参考掘金应用）
// go run main.go init 初始化用户计数值
// go run main.go get 1556564194374926  // 打印用户(1556564194374926)的所有计数值
// go run main.go incr_like 1556564194374926 // 点赞数+1
// go run main.go incr_collect 1556564194374926 // 点赞数+1
// go run main.go decr_like 1556564194374926 // 点赞数-1
// go run main.go decr_collect 1556564194374926 // 点赞数-1
func Ex05(ctx context.Context, args []string) {
	if len(args) == 0 {
		fmt.Printf("args can NOT be empty\n")
		os.Exit(1)
	}
	arg1 := args[0]
	switch arg1 {
	case "init":
		Ex06InitUserCounter(ctx)
	case "get":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		GetUserCounter(ctx, userID)
	case "incr_like":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		IncrByUserLike(ctx, userID)
	case "incr_collect":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		IncrByUserCollect(ctx, userID)
	case "decr_like":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		DecrByUserLike(ctx, userID)
	case "decr_collect":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		DecrByUserCollect(ctx, userID)
	}

}

func Ex06InitUserCounter(ctx context.Context) {
	pipe := RedisClient.Pipeline()
	userCounters := []map[string]interface{}{
		{"user_id": "1556564194374926", "got_digg_count": 10693, "got_view_count": 2238438, "followee_count": 176, "follower_count": 9895, "follow_collect_set_count": 0, "subscribe_tag_count": 95},
		{"user_id": "1111", "got_digg_count": 19, "got_view_count": 4},
		{"user_id": "2222", "got_digg_count": 1238, "follower_count": 379},
	}
	for _, counter := range userCounters {
		uid, err := strconv.ParseInt(counter["user_id"].(string), 10, 64)
		key := GetUserCounterKey(uid)
		rw, err := pipe.Del(ctx, key).Result()
		if err != nil {
			fmt.Printf("del %s, rw=%d\n", key, rw)
		}
		_, err = pipe.HMSet(ctx, key, counter).Result()
		if err != nil {
			panic(err)
		}

		fmt.Printf("设置 uid=%d, key=%s\n", uid, key)
	}
	// 批量执行上面for循环设置好的hmset命令
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func GetUserCounterKey(userID int64) string {
	return fmt.Sprintf("%s_%d", Ex05UserCountKey, userID)
}

func GetUserCounter(ctx context.Context, userID int64) {
	pipe := RedisClient.Pipeline()
	GetUserCounterKey(userID)
	pipe.HGetAll(ctx, GetUserCounterKey(userID))
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, cmder := range cmders {
		counterMap, err := cmder.(*redis.MapStringStringCmd).Result()
		if err != nil {
			panic(err)
		}
		for field, value := range counterMap {
			fmt.Printf("%s: %s\n", field, value)
		}
	}
}

// IncrByUserLike 点赞数+1
func IncrByUserLike(ctx context.Context, userID int64) {
	incrByUserField(ctx, userID, "got_digg_count")
}

// IncrByUserCollect 收藏数+1
func IncrByUserCollect(ctx context.Context, userID int64) {
	incrByUserField(ctx, userID, "follow_collect_set_count")
}

// DecrByUserLike 点赞数-1
func DecrByUserLike(ctx context.Context, userID int64) {
	decrByUserField(ctx, userID, "got_digg_count")
}

// DecrByUserCollect 收藏数-1
func DecrByUserCollect(ctx context.Context, userID int64) {
	decrByUserField(ctx, userID, "follow_collect_set_count")
}

func incrByUserField(ctx context.Context, userID int64, field string) {
	change(ctx, userID, field, 1)
}

func decrByUserField(ctx context.Context, userID int64, field string) {
	change(ctx, userID, field, -1)
}

func change(ctx context.Context, userID int64, field string, incr int64) {
	redisKey := GetUserCounterKey(userID)
	before, err := RedisClient.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	beforeInt, err := strconv.ParseInt(before, 10, 64)
	if err != nil {
		panic(err)
	}
	if beforeInt+incr < 0 {
		fmt.Printf("禁止变更计数，计数变更后小于0. %d + (%d) = %d\n", beforeInt, incr, beforeInt+incr)
		return
	}
	fmt.Printf("user_id: %d\n更新前\n%s = %s\n--------\n", userID, field, before)
	_, err = RedisClient.HIncrBy(ctx, redisKey, field, incr).Result()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("更新记录[%d]:%d\n", userID, num)
	count, err := RedisClient.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_id: %d\n更新后\n%s = %s\n--------\n", userID, field, count)
}
