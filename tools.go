package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成四位随机数
func getRandnum() string {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子

	// 生成一个范围在 0 到 9999 之间的随机数
	randomNumber := rand.Intn(10000)

	// 格式化输出为四位数，前面可以有零
	randomNumberString := fmt.Sprintf("%04d", randomNumber)
	return randomNumberString
}