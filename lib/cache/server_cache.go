/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 15:23
 */

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/link1st/gowebsocket/lib/redislib"
	"github.com/link1st/gowebsocket/models"
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
func SetServerInfo(server *models.Server, currentTime uint64) (err error) {
	key := getServersHashKey()
	value := fmt.Sprintf("%d", currentTime)
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)
		return
	}
	redisClient.Do(context.Background(), "Expire", key, serversHashCacheTime)
	return
}

// 下线服务器信息
func DelServerInfo(server *models.Server) (err error) {
	key := getServersHashKey()
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do(context.Background(), "Expire", key, serversHashCacheTime)

	return
}

func GetServerAll(currentTime uint64) (servers []*models.Server, err error) {

	servers = make([]*models.Server, 0)
	key := getServersHashKey()

	redisClient := redislib.GetClient()

	val, err := redisClient.Do(context.Background(), "hGetAll", key).Result()

	valByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valByte))

	serverMap, err := redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)

		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		// 超时
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}

		server, err := models.StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
