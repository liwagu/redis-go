package example

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

const Ex06RankKey = "ex06_rank_zset"

type Ex06ItemScore struct {
	ItemNam string
	Score   float64
}

// Ex06 排行榜
// go run main.go init // 初始化积分
// go run main.go Ex06 rev_order // 输出完整榜单
// go run main.go  Ex06 order_page 0 // 逆序分页输出，offset=1
// go run main.go  Ex06 get_rank user2 // 获取user2的排名
// go run main.go  Ex06 get_score user2 // 获取user2的分数
// go run main.go  Ex06 add_user_score user2 10 // 为user2设置为10分
// zadd ex06_rank_zset 15 andy
// zincrby ex06_rank_zset -9 andy // andy 扣9分，排名掉到最后一名
func Ex06(ctx context.Context, args []string) {
	arg1 := args[0]
	switch arg1 {
	case "init":
		Ex06InitUserScore(ctx)
	case "rev_order":
		GetRevOrderAllList(ctx, 0, -1)
	case "order_page":
		pageSize := int64(2)
		if len(args[1]) > 0 {
			offset, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				panic(err)
			}
			GetOrderListByPage(ctx, offset, pageSize)
		}
	case "get_rank":
		GetUserRankByName(ctx, args[1])
	case "get_score":
		GetUserScoreByName(ctx, args[1])
	case "add_user_score":
		if len(args) < 3 {
			fmt.Printf("参数错误，可能是缺少需要增加的分值。eg：go run main.go  Ex06 add_user_score user2 10\n")
			return
		}
		score, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			panic(err)
		}
		AddUserScore(ctx, args[1], score)
	}
	return
}

func Ex06InitUserScore(ctx context.Context) {
	initList := []redis.Z{
		{Member: "user1", Score: 10}, {Member: "user2", Score: 232}, {Member: "user3", Score: 129},
		{Member: "user4", Score: 232},
	}
	// 清空榜单
	if err := RedisClient.Del(ctx, Ex06RankKey).Err(); err != nil {
		panic(err)
	}

	nums, err := RedisClient.ZAdd(ctx, Ex06RankKey, initList...).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("初始化榜单Item数量:%d\n", nums)
}

// 榜单逆序输出
// ZRANGE ex06_rank_zset +inf -inf BYSCORE  rev WITHSCORES
// 正序输出
// ZRANGE ex06_rank_zset 0 -1 WITHSCORES
func GetRevOrderAllList(ctx context.Context, limit, offset int64) {
	resList, err := RedisClient.ZRevRangeWithScores(ctx, Ex06RankKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n榜单:\n")
	for i, z := range resList {
		fmt.Printf("第%d名 %s\t%f\n", i+1, z.Member, z.Score)
	}
}

func GetOrderListByPage(ctx context.Context, offset, pageSize int64) {
	// zrange ex06_rank_zset 300 0 byscore rev limit 1 2 withscores // 取300分到0分之间的排名
	// zrange ex06_rank_zset -inf +inf byscore withscores 正序输出
	// ZRANGE ex06_rank_zset +inf -inf BYSCORE  REV WITHSCORES 逆序输出所有排名
	// zrange ex06_rank_zset +inf -inf byscore rev limit 0 2 withscores 逆序分页输出排名
	zRangeArgs := redis.ZRangeArgs{
		Key:     Ex06RankKey,
		ByScore: true,
		Rev:     true,
		Start:   "-inf",
		Stop:    "+inf",
		Offset:  offset,
		Count:   pageSize,
	}
	resList, err := RedisClient.ZRangeArgsWithScores(ctx, zRangeArgs).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n榜单(offest=%d, pageSize=%d):\n", offset, pageSize)
	offNum := int(pageSize * offset)
	for i, z := range resList {
		rank := i + 1 + offNum
		fmt.Printf("第%d名 %s\t%f\n", rank, z.Member, z.Score)
	}
	fmt.Println()
}

// GetUserRankByName 获取用户排名
func GetUserRankByName(ctx context.Context, name string) {
	rank, err := RedisClient.ZRevRank(ctx, Ex06RankKey, name).Result()
	if err != nil {
		fmt.Errorf("error getting name=%s, err=%v", name, err)
		return
	}
	fmt.Printf("name=%s, 排名=%d\n", name, rank+1)
}

// GetUserScoreByName 获取用户分值
func GetUserScoreByName(ctx context.Context, name string) {
	score, err := RedisClient.ZScore(ctx, Ex06RankKey, name).Result()
	if err != nil {
		fmt.Errorf("error getting name=%s, err=%v", name, err)
		return
	}
	fmt.Println(time.Now().UnixMilli())
	fmt.Printf("name=%s, 分数=%f\n", name, score)
}

// AddUserScore 排名用户
func AddUserScore(ctx context.Context, name string, score float64) {
	num, err := RedisClient.ZIncrBy(ctx, Ex06RankKey, score, name).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, add_score=%f, score=%f\n", name, score, num)
}
