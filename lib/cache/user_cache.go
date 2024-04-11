// Package cache 缓存
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/link1st/gowebsocket/v2/lib/redislib"
	"github.com/link1st/gowebsocket/v2/models"

	"github.com/redis/go-redis/v9"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
	return
}

// GetUserOnlineInfo 获取用户在线信息
func GetUserOnlineInfo(userKey string) (userOnline *models.UserOnline, err error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)
	data, err := redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("GetUserOnlineInfo", userKey, err)
			return
		}
		fmt.Println("GetUserOnlineInfo", userKey, err)
		return
	}
	userOnline = &models.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		fmt.Println("获取用户在线数据 json Unmarshal", userKey, err)
		return
	}
	fmt.Println("获取用户在线数据", userKey, "time", userOnline.LoginTime, userOnline.HeartbeatTime, "AccIp",
		userOnline.AccIp, userOnline.IsLogoff)
	return
}

// SetUserOnlineInfo 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *models.UserOnline) (err error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)
	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		fmt.Println("设置用户在线数据 json Marshal", key, err)
		return
	}
	_, err = redisClient.Do(context.Background(), "setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置用户在线数据 ", key, err)
		return
	}
	return
}
