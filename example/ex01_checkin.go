package example

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

var ctx = context.Background()

const continuesCheckKey = "cc_uid_%d"

// Ex01 连续签到天数
func Ex01(ctx context.Context, params []string) {
	if userID, err := strconv.ParseInt(params[0], 10, 64); err == nil {
		addContinuesDays(ctx, userID)
	} else {
		fmt.Printf("参数错误, params=%v, error: %v\n", params, err)
	}
}

// addContinuesDays 为用户签到续期
func addContinuesDays(ctx context.Context, userID int64) {
	key := fmt.Sprintf(continuesCheckKey, userID)
	// 1. 连续签到数+1
	err := RedisClient.Incr(ctx, key).Err()
	if err != nil {
		fmt.Errorf("用户[%d]连续签到失败", userID)
	} else {
		expAt := beginningOfDay().Add(48 * time.Hour)
		// 2. 设置签到记录在后天的0点到期
		if err := RedisClient.ExpireAt(ctx, key, expAt).Err(); err != nil {
			panic(err)
		} else {
			// 3. 打印用户续签后的连续签到天数
			day, err := getUserCheckInDays(ctx, userID)
			if err != nil {
				panic(err)
			}
			fmt.Printf("用户[%d]连续签到：%d(天), 过期时间:%s", userID, day, expAt.Format("2006-01-02 15:04:05"))
		}
	}
}

// getUserCheckInDays 获取用户连续签到天数
func getUserCheckInDays(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf(continuesCheckKey, userID)
	days, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if daysInt, err := strconv.ParseInt(days, 10, 64); err != nil {
		panic(err)
	} else {
		return daysInt, nil
	}
}

// beginningOfDay 获取今天0点时间
func beginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}
