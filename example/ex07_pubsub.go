package example

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
)

const ex07Channel = "es_ch"

// Ex04Params Ex04的自定义函数
type Ex07Params struct {
}

func Ex07(ctx context.Context) {
	pubSub := RedisClient.Subscribe(ctx, ex07Channel)
	defer func(mPubSub *redis.PubSub) {
		mPubSub.Unsubscribe(ctx, ex07Channel)
		mPubSub.Close()
	}(pubSub)
	for msg := range pubSub.Channel() {
		// 打印收到的消息
		// fmt.Println(msg.Channel)
		articleID, err := strconv.ParseInt(msg.Payload, 10, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("读取文章[%d]标题、正文，发送到ES更新索引\n", articleID)
	}
}
