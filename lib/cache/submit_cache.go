// Package cache 缓存
package cache

import (
	"context"
	"fmt"

	"github.com/link1st/gowebsocket/v2/lib/redislib"
)

const (
	submitAgainPrefix = "acc:submit:again:" // 数据不重复提交
)

// getSubmitAgainKey 获取数据提交去除key
func getSubmitAgainKey(from string, value string) (key string) {
	key = fmt.Sprintf("%s%s:%s", submitAgainPrefix, from, value)

	return
}

// submitAgain 重复提交
// return true:重复提交 false:第一次提交
func submitAgain(from string, second int, value string) (isSubmitAgain bool) {

	// 默认重复提交
	isSubmitAgain = true
	key := getSubmitAgainKey(from, value)
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "setNx", key, "1").Int()
	if err != nil {
		fmt.Println("submitAgain", key, number, err)
		return
	}
	if number != 1 {
		return
	}

	// 第一次提交
	isSubmitAgain = false
	redisClient.Do(context.Background(), "Expire", key, second)
	return

}

// SeqDuplicates Seq 重复提交
func SeqDuplicates(seq string) (result bool) {
	result = submitAgain("seq", 12*60*60, seq)
	return
}
