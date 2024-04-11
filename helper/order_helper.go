// Package helper 帮助函数
package helper

import (
	"fmt"
	"time"
)

// GetOrderIdTime 获取订单 ID
func GetOrderIdTime() (orderId string) {
	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)
	return
}
