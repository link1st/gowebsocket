/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package common

type JsonResult struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(code uint32, message string, data interface{}) JsonResult {

	message = GetErrorMessage(code, message)
	jsonMap := grantMap(code, message, data)

	return jsonMap
}

// 按照接口格式生成原数据数组
func grantMap(code uint32, message string, data interface{}) JsonResult {

	jsonMap := JsonResult{
		Code: code,
		Msg:  message,
		Data: data,
	}
	return jsonMap
}
