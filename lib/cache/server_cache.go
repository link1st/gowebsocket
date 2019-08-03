/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 15:23
 */

package cache

import (
	"fmt"
	"go-common/library/cache/redis"
	"gowebsocket/lib/redislib"
)

const (
	serversHashKey       = "acc:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60        // key过期时间
	serversHashTimeout   = 3 * 60             // 超时时间
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)

	return
}

// 设置服务器信息
func SetServerInfo(field string, currentTime uint64) (err error) {
	key := getServersHashKey()

	value := fmt.Sprintf("%d", currentTime)

	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hSet", key, field, value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

// 下线服务器信息
func DelServerInfo(field string) (err error) {
	key := getServersHashKey()
	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hDel", key, field).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

func GetServerAll(currentTime uint64) (servers []string, err error) {

	servers = make([]string, 0)
	key := getServersHashKey()

	redisClient := redislib.GetClient()
	serverMap, err := redis.IntMap(redisClient.Do("hGetAll", key).Result())
	if err != nil {
		fmt.Println("SetServerInfo", key, err)

		return
	}

	for key, value := range serverMap {
		// 超时
		if uint64(value)+serversHashTimeout <= currentTime {
			continue
		}
		servers = append(servers, key)
	}

	return
}
