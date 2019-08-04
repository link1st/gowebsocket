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

func GetOrderIdTime() (orderId string) {

	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)

	return
}
