/**
* Created by GoLand.
* User: link1st
* Date: 2021/2/21
* Time: 12:45
 */

package helper

import "strconv"

// StrToUint32 str to uint32
func StrToUint32(str string) uint32 {
	data, _ := strconv.ParseInt(str, 10, 32)
	return uint32(data)
}
