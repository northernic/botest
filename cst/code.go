package cst

import "time"

var (
	StartTime        = time.Time{}        // 每天的开始时间（小时）
	EndTime          = time.Time{}        // 次日的结束时间（小时）
	SendCodeInterval = 3600 * time.Second // 发送兑换码的间隔时间

	CodeNum int64 = 1     // 每次发送的兑换码数量
	IsSend        = false // 是否发送兑换码 默认不发送

	ActivityText = "" // 活动文案

	SplitSigns = []string{",", ";", ":", "|", "-", "\n", " ", "\t"} // 兑换码分割符
)

func GetTicker() *time.Ticker {
	return time.NewTicker(SendCodeInterval)
}
