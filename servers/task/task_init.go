/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-29
* Time: 14:15
 */

package task

import "time"

type TimerFunc func(interface{}) bool

/**
 * 定时调用
 * @delay 首次延时
 * @tick  间隔
 * @fun   定时执行function
 * @param fun参数
 */
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}, funcDefer TimerFunc, paramDefer interface{}) {
	go func() {
		defer func() {
			if funcDefer != nil {
				funcDefer(paramDefer)
			}
		}()

		if fun == nil {
			return
		}

		t := time.NewTimer(delay)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				if fun(param) == false {
					return
				}
				t.Reset(tick)
			}
		}
	}()
}
