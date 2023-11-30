package cst

import "time"

const (
	BeijingOffset = 8 * time.Hour //北京
	BrazilOffset  = 3 * time.Hour //巴西
)

// TodayFormatStrOffset 今天日期字符串
func TodayFormatStrOffset(offset time.Duration) string {
	now := time.Now().Add(offset)
	return now.Format(time.DateOnly)
}
