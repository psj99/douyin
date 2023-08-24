package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	UserInfoCountKey = "user_info_count" // key前缀
	FollowCount      = "followee_count"  //关注数
	FanCount         = "follower_count"  //粉丝数
	GotNumbsCount    = "got_digg_count"  //获赞数
)

func getUserCounterKey(userID int64) string {
	return fmt.Sprintf("%s_%d", UserInfoCountKey, userID)
}

func InitUserCounter(ctx context.Context) {
	pipe := _redis.Pipeline()
	userCounters := []map[string]interface{}{
		{"user_id": "1556564194374926", "got_digg_count": 10693, "follower_count": 9895},
		{"user_id": "1111", "got_digg_count": 19},
		{"user_id": "2222", "got_digg_count": 1238},
	}
	for _, counter := range userCounters {
		uid, _ := strconv.ParseInt(counter["user_id"].(string), 10, 64)
		key := getUserCounterKey(uid)
		rw, err := pipe.Del(ctx, key).Result()
		if err != nil {
			fmt.Printf("del %s, rw=%d\n", key, rw)
		}
		_, err = pipe.HMSet(ctx, key, counter).Result()
		if err != nil {
			panic(err)
		}
		fmt.Printf("设置 uid=%d, key=%s\n", uid, key)
	}
	// 批量执行上面for循环设置好的hmset命令
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}

// GetUserCounter 读取UserID的所有键值信息
func GetUserCounter(ctx context.Context, userID int64) {
	pipe := _redis.Pipeline()
	pipe.HGetAll(ctx, getUserCounterKey(userID))
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("用户不存在")
			return
		}
		panic(err)
	}
	for _, cmder := range cmders {
		counterMap, err := cmder.(*redis.MapStringStringCmd).Result()
		if err != nil {
			panic(err)
		}
		for field, value := range counterMap {
			fmt.Printf("%s: %s\n", field, value)
		}
	}
}

// IncrByUserLike 点赞数+1GetUserCounter
func IncrByUserLike(ctx context.Context, userID int64) {
	incrByUserField(ctx, userID, GotNumbsCount)
}

// DecrByUserLike 点赞数-1
func DecrByUserLike(ctx context.Context, userID int64) {
	decrByUserField(ctx, userID, GotNumbsCount)
}

// GetByUserLike 获取用户获赞数
func GetByUserLike(ctx context.Context, userID int64) string {
	return getByUserFiled(ctx, userID, GotNumbsCount)
}

// IncrByUserFan 粉丝数+1
func IncrByUserFan(ctx context.Context, userID int64) {
	incrByUserField(ctx, userID, FanCount)
}

// DecrByUserFan  粉丝数-1
func DecrByUserFan(ctx context.Context, userID int64) {
	decrByUserField(ctx, userID, FanCount)
}

// GetByUserFan 获取用户粉丝数
func GetByUserFan(ctx context.Context, userID int64) string {
	return getByUserFiled(ctx, userID, FanCount)
}

// IncryByUserFolllow 关注数+1
func IncryByUserFolllow(ctx context.Context, userID int64) {
	incrByUserField(ctx, userID, FollowCount)
}

// DecrByUserFolllow 关注数-1
func DecrByUserFolllow(ctx context.Context, userID int64) {
	decrByUserField(ctx, userID, FollowCount)
}

// GetByUserFollow 获取关注数
func GetByUserFollow(ctx context.Context, userID int64) string {
	return getByUserFiled(ctx, userID, FollowCount)
}

func incrByUserField(ctx context.Context, userID int64, field string) {
	change(ctx, userID, field, 1)
}

func decrByUserField(ctx context.Context, userID int64, field string) {
	change(ctx, userID, field, -1)
}

func getByUserFiled(ctx context.Context, userID int64, field string) string {
	redisKey := getUserCounterKey(userID)
	count, err := _redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			count = "0"
		} else {
			panic(err)
		}
	}
	return count
}

func change(ctx context.Context, userID int64, field string, incr int64) {
	redisKey := getUserCounterKey(userID)
	before, err := _redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	beforeInt, err := strconv.ParseInt(before, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	if beforeInt+incr < 0 {
		fmt.Printf("禁止变更计数，计数变更后小于0. %d + (%d) = %d\n", beforeInt, incr, beforeInt+incr)
		return
	}
	fmt.Printf("user_id: %d\n更新前\n%s = %s\n--------\n", userID, field, before)
	_, err = _redis.HIncrBy(ctx, redisKey, field, incr).Result()
	if err != nil {
		panic(err)
	}
	count, err := _redis.HGet(ctx, redisKey, field).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("user_id: %d\n更新后\n%s = %s\n--------\n", userID, field, count)
}
