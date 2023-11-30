package cst

import (
	"regexp"
	"time"
)

var (
	StartTime               = time.Time{} // 每天的开始时间（小时）
	EndTime                 = time.Time{} // 次日的结束时间（小时）
	SendCodeInterval string = "1h"        // 发送兑换码的间隔时间

	CodeNum int64 = 1     // 每次发送的兑换码数量
	IsSend        = false // 是否发送兑换码 默认不发送

	ActivityText = "" // 活动文案

	SplitSigns = []string{",", ";", ":", "|", "-", "\n", " ", "\t"} // 兑换码分割符
)

func CheckInterval(s string) bool {
	re := regexp.MustCompile(`^(\d+h)?(\d+m)?(\d+s)?$`)
	return re.MatchString(s)
}
