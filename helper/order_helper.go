/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-01
* Time: 18:13
 */

package helper

import (
	"fmt"
	"time"
)

// GetOrderIDTime 生成订单号
func GetOrderIDTime() (orderId string) {
	currentTime := time.Now().Nanosecond()
	return fmt.Sprintf("%d", currentTime)
}
