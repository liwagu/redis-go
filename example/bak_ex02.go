package example

//
// // setnx 分布式锁
// import (
// 	"context"
// 	"fmt"
// 	"sort"
// 	"strconv"
// 	"sync"
// 	"time"
// )
//
// const resourceKey = "syncKey"      // 分布式锁的key
// const exp = 800 * time.Millisecond // 锁的过期时间，避免死锁
//
// // EventLog 搜集日志的结构
// type EventLog struct {
// 	eventTime time.Time
// 	log       string
// }
//
// var eventLogs []EventLog
//
// // Ex02 只是体验SetNX的特性，不是高可用的分布式锁实现
// // 该实现存在的问题:
// // (1) 业务超时解锁，导致并发问题。业务执行时间超过锁超时时间
// // (2) redis主备切换临界点问题。主备切换后，A持有的锁还未同步到新的主节点时，B可在新主节点获取锁，导致并发问题。
// // (3) redis集群脑裂，导致出现多个主节点
// func Ex02(ctx context.Context) {
// 	wg := &sync.WaitGroup{}
// 	rts := 10
// 	eventLogs = make([]EventLog, 0, rts)
// 	for i := 0; i < rts; i++ {
// 		wg.Add(1)
// 		// 启动协程模拟并发逻辑
// 		go ex02Work(ctx, i, wg)
// 	}
// 	wg.Wait()
//
// 	// 日志按时间正序输出
// 	sort.Slice(eventLogs, func(i, j int) bool {
// 		return eventLogs[i].eventTime.Before(eventLogs[j].eventTime)
// 	})
// 	for i := range eventLogs {
// 		fmt.Println(eventLogs[i].log)
// 		if (i+1)%3 == 0 {
// 			fmt.Println()
// 		}
// 	}
// }
//
// func ex02Work(ctx context.Context, routine int, wg *sync.WaitGroup) {
// 	defer ex02ReleaseLock(ctx, routine, wg) // 4. 释放锁
// 	for {
// 		// 1. 尝试获取锁
// 		// exp - 锁过期设置,避免异常死锁
// 		acquired, err := RedisClient.SetNX(ctx, resourceKey, routine, exp).Result() // 尝试获取锁
// 		if err != nil {
// 			eventLogs = append(eventLogs, EventLog{
// 				time.Now(), fmt.Sprintf("[%s] error routine[%d], %v", time.Now().Format(time.RFC3339Nano), routine, err),
// 			})
// 			panic(err)
// 		}
// 		if acquired {
// 			// 2. 成功获取锁
// 			eventLogs = append(eventLogs, EventLog{
// 				time.Now(), fmt.Sprintf("[%s] routine[%d] 获取锁", time.Now().Format(time.RFC3339Nano), routine),
// 			})
// 			// 3. sleep 模拟业务逻辑耗时
// 			time.Sleep(10 * time.Millisecond)
// 			eventLogs = append(eventLogs, EventLog{
// 				time.Now(), fmt.Sprintf("[%s] routine[%d] 完成业务逻辑", time.Now().Format(time.RFC3339Nano), routine),
// 			})
// 			return
// 		} else {
// 			// 没有获得锁，等待后重试
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}
// }
//
// func ex02ReleaseLock(ctx context.Context, routine int, wg *sync.WaitGroup) {
// 	routineMark, _ := RedisClient.Get(ctx, resourceKey).Result()
// 	if strconv.FormatInt(int64(routine), 10) != routineMark {
// 		// 其它协程误删lock
// 		panic(fmt.Sprintf("del err lock[%s] can not del by [%d]", routineMark, routine))
// 	}
// 	set, err := RedisClient.Del(ctx, resourceKey).Result()
// 	if set == 1 {
// 		eventLogs = append(eventLogs, EventLog{
// 			time.Now(), fmt.Sprintf("[%s] routine[%d] 释放锁", time.Now().Format(time.RFC3339Nano), routine),
// 		})
// 	} else {
// 		eventLogs = append(eventLogs, EventLog{
// 			time.Now(), fmt.Sprintf("[%s] routine[%d] no lock to del", time.Now().Format(time.RFC3339Nano), routine),
// 		})
// 	}
// 	if err != nil {
// 		fmt.Errorf("[%s] error routine=%d, %v", time.Now().Format(time.RFC3339Nano), routine, err)
// 		panic(err)
// 	}
// 	wg.Done()
// }
