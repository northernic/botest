package cst

import "time"

var StartTime = time.Time{} // 每天的开始时间（小时）
var EndTime = time.Time{}   // 次日的结束时间（小时）
var IsSend = false          // 是否发送兑换码 默认不发送

var ActivityText = ""

func GetTicker() *time.Ticker {
	return time.NewTicker(15 * time.Second)
}
