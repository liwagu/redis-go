package common

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type ConcurrentEventLogger struct {
	eventLogs []EventLog
}

// EventLog 搜集日志的结构
type EventLog struct {
	EventTime time.Time
	Log       string
}

func NewConcurrentEventLog(ctx context.Context, logsLength int) *ConcurrentEventLogger {
	if logsLength <= 0 {
		logsLength = 32
	}
	logContainer := make([]EventLog, 0, logsLength)
	return &ConcurrentEventLogger{eventLogs: logContainer}
}

// Append 追加日志
func (ceLog *ConcurrentEventLogger) Append(mLog EventLog) {
	ceLog.eventLogs = append(ceLog.eventLogs, mLog)
}

// PrintLogs  日志按时间正序输出
func (ceLog *ConcurrentEventLogger) PrintLogs() {
	sort.Slice(ceLog.eventLogs, func(i, j int) bool {
		return ceLog.eventLogs[i].EventTime.Before(ceLog.eventLogs[j].EventTime)
	})
	for i := range ceLog.eventLogs {
		fmt.Println(ceLog.eventLogs[i].Log)
		// if (i+1)%3 == 0 {
		// 	fmt.Println()
		// }
	}
}

// LogFormat 包含通用日志前缀 [2022-11-27T12:36:00.213454+08:00] routine[5]
func LogFormat(routine int, format string, a ...any) string {
	tpl := "[%s] routine[%d] " + format
	sr := []any{time.Now().Format(time.RFC3339Nano), routine}
	sr = append(sr, a...)
	return fmt.Sprintf(tpl, sr...)
}
