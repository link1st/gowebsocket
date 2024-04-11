// Package helper 帮助函数
package helper

import (
	"fmt"
	"time"
)

// GetOrderIDTime 获取订单 ID
func GetOrderIDTime() (orderID string) {
	currentTime := time.Now().Nanosecond()
	orderID = fmt.Sprintf("%d", currentTime)
	return
}
