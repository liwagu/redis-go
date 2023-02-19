package example

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDel(t *testing.T) {
	Convey("Testing", t, func() {
		ctx := context.Background()
		set, err := RedisClient.Del(ctx, "pky").Result()
		So(err, ShouldBeNil)
		So(set, ShouldEqual, 1)
	})
}

func TestLjkl(t *testing.T) {
	ctx := context.Background()
	Convey("fff", t, func() {
		// 当jjly未初始化时
		result, err := RedisClient.Get(ctx, "jjly").Result()
		So(err, ShouldEqual, redis.Nil)
		So(result, ShouldEqual, "")
	})
}
