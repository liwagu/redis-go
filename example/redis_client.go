package example

import (
	"github.com/go-redis/redis/v9"
)

var RedisClient *redis.Client

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	RedisClient = rdb
}
