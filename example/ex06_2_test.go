package example

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 命令行运行测试用例
// (1) cd redis_course
// (2) go test -v -run TestSetScoreWithTime example/ex06_2_test.go example/ex06_2.go example/redis_client.go
func TestSetScoreWithTime(t *testing.T) {
	Convey("测试-时间信息写入与获取", t, func() {
		ctx := context.Background()
		score := (1 << 21) - 1          // 2,097,151
		So(score, ShouldEqual, 2097151) // 支持最大的分值
		scoreWithTime := GetScoreWithTime(ctx, int64(score))
		scoreInt64 := GetScoreWithoutTime(ctx, scoreWithTime)
		So(scoreInt64, ShouldEqual, int64(2097151))
	})
}
