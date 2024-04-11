// Package redislib redis 库
package redislib

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
)

// NewClient 初始化 Redis 客户端
func NewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	pong, err := client.Ping(context.Background()).Result()
	fmt.Println("初始化redis:", pong, err)
}

// GetClient 获取客户端
func GetClient() (c *redis.Client) {
	return client
}
