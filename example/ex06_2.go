package example

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

// 0（最高位不用）|
// __000000 00000000 00000000（22bit表示分值），最大有效值 2,097,151
// 0 00000000 00000000 00000000 00000000 00000000（41bit表示毫秒级时间戳，毫秒级时间戳占用41bit）

const Ex062RankKey = "ex062_rank_zset"
const MaxTime = (1 << 41) - 1 // 1 11111111 11111111 11111111 11111111 11111111

// GetScoreWithTime 获取携带时间信息的分值
func GetScoreWithTime(ctx context.Context, score int64) int64 {
	scoreWithTime := score << 41
	now := time.Now().UnixMilli()
	scoreWithTime |= MaxTime - now // 写入时间约晚，(MaxTime - now)越小，使得后写入的数据低位越小。相同分值情况下，后写入的排序就越靠后
	return scoreWithTime
}

// GetScoreWithoutTime 去除数值中的时间信息
func GetScoreWithoutTime(ctx context.Context, scoreWithTime int64) int64 {
	score := scoreWithTime >> 41
	return score
}

// Ex06_2 排行榜，支持同分值按写入时间序排序
// go run main.go init // 初始化积分
// go run main.go rev_order // 输出完整榜单
// go run main.go order_page // 逆序分页输出
// go run main.go order_page 1 // 逆序分页输出，offset=1
// go run main.go add_user_score judy 23
// zadd ex062_rank_zset 15 andy
// zincrby ex062_rank_zset -9 andy // andy 扣9分，排名掉到最后一名
func Ex06_2(ctx context.Context, args []string) {
	arg1 := args[0]
	switch arg1 {
	case "init":
		Ex062InitUserScore(ctx)
	case "rev_order":
		GetRevOrderAllList062(ctx, 0, -1)
	case "order_page":
		pageSize := int64(2)
		offset := int64(0)
		var err error
		if len(args[1]) > 0 {
			offset, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				panic(err)
			}
		}
		GetOrderListByPage062(ctx, offset, pageSize)
	case "get_rank":
		GetUserRankByName062(ctx, args[1])
	case "get_score":
		GetUserScoreByName062(ctx, args[1])
	case "add_user_score":
		score, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			panic(err)
		}
		AddUserScore062(ctx, args[1], score)
	}
	return
}

// Ex062InitUserScore 初始化榜单
func Ex062InitUserScore(ctx context.Context) {
	initList := []redis.Z{
		{Member: "user1", Score: float64(GetScoreWithTime(ctx, 10))},
		{Member: "user2", Score: float64(GetScoreWithTime(ctx, 232))},
		{Member: "user3", Score: float64(GetScoreWithTime(ctx, 129))},
	}
	// 清空榜单
	if err := RedisClient.Del(ctx, Ex062RankKey).Err(); err != nil {
		panic(err)
	}

	nums, err := RedisClient.ZAdd(ctx, Ex062RankKey, initList...).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("初始化榜单Item数量:%d\n", nums)
	time.Sleep(2 * time.Second)
	lastItem := []redis.Z{
		{Member: "user4", Score: float64(GetScoreWithTime(ctx, 232))},
	}
	nums, err = RedisClient.ZAdd(ctx, Ex062RankKey, lastItem...).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("初始化榜单Item数量:%d\n", nums)
}

// 榜单逆序输出
// ZRANGE ex06_rank_zset +inf -inf BYSCORE  rev WITHSCORES
// 正序输出
// ZRANGE ex06_rank_zset 0 -1 WITHSCORES
func GetRevOrderAllList062(ctx context.Context, limit, offset int64) {
	resList, err := RedisClient.ZRevRangeWithScores(ctx, Ex062RankKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n榜单:\n")
	for i, z := range resList {
		scoreWithoutTime := GetScoreWithoutTime(ctx, int64(z.Score))
		fmt.Printf("第%d名 %s\t%d\n", i+1, z.Member, scoreWithoutTime)
	}
}

func GetOrderListByPage062(ctx context.Context, offset, pageSize int64) {
	// zrange ex06_rank_zset 300 0 byscore rev limit 1 2 withscores // 取300分到0分之间的排名
	// zrange ex06_rank_zset -inf +inf byscore withscores 正序输出
	// ZRANGE ex06_rank_zset +inf -inf BYSCORE  REV WITHSCORES 逆序输出所有排名
	// zrange ex06_rank_zset +inf -inf byscore rev limit 0 2 withscores 逆序分页输出排名
	zRangeArgs := redis.ZRangeArgs{
		Key:     Ex062RankKey,
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
		scoreWithoutTime := GetScoreWithoutTime(ctx, int64(z.Score))
		fmt.Printf("第%d名 %s\t%d\n", rank, z.Member, scoreWithoutTime)
	}
	fmt.Println()
}

// GetUserRankByName062 获取用户排名
func GetUserRankByName062(ctx context.Context, name string) {
	rank, err := RedisClient.ZRevRank(ctx, Ex062RankKey, name).Result()
	if err != nil {
		fmt.Errorf("error getting name=%s, err=%v", name, err)
		return
	}
	fmt.Printf("name=%s, 排名=%d\n", name, rank+1)
}

// GetUserScoreByName062 获取用户分值
func GetUserScoreByName062(ctx context.Context, name string) {
	// fmt.Println(strconv.FormatInt(1<<41, 2))
	// fmt.Println(strconv.FormatInt(MaxTime, 2))
	score, err := RedisClient.ZScore(ctx, Ex062RankKey, name).Result()
	if err != nil {
		fmt.Errorf("error getting name=%s, err=%v", name, err)
		return
	}
	scoreWithoutTime := GetScoreWithoutTime(ctx, int64(score))
	fmt.Printf("name=%s, 分数=%d\n", name, scoreWithoutTime)
}

// AddUserScore062 增加一个排名用户
func AddUserScore062(ctx context.Context, name string, score int64) {
	scoreWithT := GetScoreWithTime(ctx, score)
	z := redis.Z{
		Score:  float64(scoreWithT),
		Member: name,
	}
	num, err := RedisClient.ZAdd(ctx, Ex062RankKey, z).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("set[%d] name=%s, score=%d, scoreWithTime=%d\n", num, name, score, scoreWithT)
}
