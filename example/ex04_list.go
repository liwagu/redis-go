package example

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitee.com/wedone/redis_course/example/common"
)

const ex04ListenList = "ex04_list_0" // lpush ex04_list_0 AA BB

// Ex04Params Ex04的自定义函数
type Ex04Params struct {
}

func Ex04(ctx context.Context) {
	eventLogger := &common.ConcurrentEventLogger{}
	// new一个并发执行器
	// routineNums是消费端的数量，多消费的场景，可以使用ex04ConsumerPop，使用ex04ConsumerRange存在消息重复消费的问题。
	cInst := common.NewConcurrentRoutine(1, eventLogger)
	// 并发执行用户自定义函数work
	cInst.Run(ctx, Ex04Params{}, ex04ConsumerPop)
	// 按日志时间正序打印日志
	eventLogger.PrintLogs()
}

// ex04ConsumerPop 使用rpop逐条消费队列中的信息，数据从队列中移除
// 生成端使用：lpush ex04_list_0 AA BB
func ex04ConsumerPop(ctx context.Context, cInstParam common.CInstParams) {
	routine := cInstParam.Routine
	for {
		items, err := RedisClient.BRPop(ctx, 0, ex04ListenList).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(common.LogFormat(routine, "读取文章[%s]标题、正文，发送到ES更新索引", items[1]))
		// 将文章内容推送到ES
		time.Sleep(1 * time.Second)
	}
}

// ex04ConsumerRange 使用lrange批量消费队列中的数据，数据保留在队列中
// 生成端使用：rpush ex04_list_0 AA BB
// 消费端：
// 方法1 lrange ex04_list_0 -3 -1 // 从FIFO队尾中一次消费3条信息
// 方法2 rpop ex04_list_0 3
func ex04ConsumerRange(ctx context.Context, cInstParam common.CInstParams) {
	routine := cInstParam.Routine
	consumeBatchSize := int64(3) // 一次取N个消息
	for {
		// 从index(-consumeBatchSize)开始取，直到最后一个元素index(-1)
		items, err := RedisClient.LRange(ctx, ex04ListenList, -consumeBatchSize, -1).Result()
		if err != nil {
			panic(err)
		}
		if len(items) > 0 {
			fmt.Println(common.LogFormat(routine, "收到信息:%s", strings.Join(items, "->")))
			// 清除已消费的队列
			// 方法1 使用LTrim
			// 保留从index(0)开始到index(-(consumeBatchSize + 1))的部分，即为未消费的部分
			// RedisClient.LTrim(ctx, ex04ListenList, 0, -(consumeBatchSize + 1))

			// 方法2 使用RPop
			RedisClient.RPopCount(ctx, ex04ListenList, int(consumeBatchSize))
		}
		time.Sleep(3 * time.Second)
	}
}
