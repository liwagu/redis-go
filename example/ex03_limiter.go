package example

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v9"

	"gitee.com/wedone/redis_course/example/common"
)

type Ex03Params struct {
}

var ex03LimitKeyPrefix = "comment_freq_limit"
var accessQueryNum = int32(0)

const ex03MaxQPS = 10 // 限流次数

// ex03LimitKey 返回key格式为：comment_freq_limit-1669524458 // 用来记录这1秒内的请求数量
func ex03LimitKey(currentTimeStamp time.Time) string {
	return fmt.Sprintf("%s-%d", ex03LimitKeyPrefix, currentTimeStamp.Unix())
}

// Ex03 简单限流
func Ex03(ctx context.Context) {
	eventLogger := &common.ConcurrentEventLogger{}
	// new一个并发执行器
	cInst := common.NewConcurrentRoutine(500, eventLogger)
	// 并发执行用户自定义函数work
	cInst.Run(ctx, Ex03Params{}, ex03Work)
	// 按日志时间正序打印日志
	eventLogger.PrintLogs()
	fmt.Printf("放行总数：%d\n", accessQueryNum)

	fmt.Printf("\n------\n下一秒请求\n------\n")
	accessQueryNum = 0
	time.Sleep(1 * time.Second)
	// new一个并发执行器
	cInst = common.NewConcurrentRoutine(10, eventLogger)
	// 并发执行用户自定义函数work
	cInst.Run(ctx, Ex03Params{}, ex03Work)
	// 按日志时间正序打印日志
	eventLogger.PrintLogs()
	fmt.Printf("放行总数：%d\n", accessQueryNum)
}

func ex03Work(ctx context.Context, cInstParam common.CInstParams) {
	routine := cInstParam.Routine
	eventLogger := cInstParam.ConcurrentEventLogger
	key := ex03LimitKey(time.Now())
	currentQPS, err := RedisClient.Incr(ctx, key).Result()
	if err != nil || err == redis.Nil {
		err = RedisClient.Incr(ctx, ex03LimitKey(time.Now())).Err()
		if err != nil {
			panic(err)
		}
	}
	if currentQPS > ex03MaxQPS {
		// 超过流量限制，请求被限制
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       common.LogFormat(routine, "被限流[%d]", currentQPS),
		})
		// sleep 模拟业务逻辑耗时
		time.Sleep(50 * time.Millisecond)
		err = RedisClient.Decr(ctx, key).Err()
		if err != nil {
			panic(err)
		}
	} else {
		// 流量放行
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       common.LogFormat(routine, "流量放行[%d]", currentQPS),
		})
		atomic.AddInt32(&accessQueryNum, 1)
		time.Sleep(20 * time.Millisecond)
	}
}
