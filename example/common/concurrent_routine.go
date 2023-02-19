package common

import (
	"context"
	"sync"
)

// ConcurrentRoutine 并发执行器对象定义
type ConcurrentRoutine struct {
	routineNums           int                    // 定义并发协程的数量
	concurrentEventLogger *ConcurrentEventLogger // 并发日志搜集器
}

// CInstParams 定义传入callBack的参数
type CInstParams struct {
	Routine               int // 协程编号
	ConcurrentEventLogger *ConcurrentEventLogger
	CustomParams          interface{} // 用户自定义参数
}

type callBack func(ctx context.Context, params CInstParams) // 定义一个用户自定义执行函数

// NewConcurrentRoutine 初始化一个并发执行器
func NewConcurrentRoutine(routineNums int, concurrentEventLog *ConcurrentEventLogger) *ConcurrentRoutine {
	return &ConcurrentRoutine{
		routineNums: routineNums, concurrentEventLogger: concurrentEventLog,
	}
}

// Run 并发执行用户自定义函数 workFun
func (cInst *ConcurrentRoutine) Run(ctx context.Context, customParams interface{}, workFun callBack) {
	wg := &sync.WaitGroup{}
	for i := 0; i < cInst.routineNums; i++ {
		wg.Add(1)
		// 启动协程模拟并发逻辑
		go func(mCtx context.Context, mRoutine int, mParams interface{}) {
			defer wg.Done()
			workFun(mCtx, CInstParams{Routine: mRoutine, ConcurrentEventLogger: cInst.concurrentEventLogger, CustomParams: mParams})
		}(ctx, i, customParams)
	}
	wg.Wait()
}
