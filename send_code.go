package main

import (
	"bot/cst"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// 定义时间范围的起止时间

func sendCodeTimer(bot *tgbotapi.BotAPI) {

	//next := getNextExecutionTime()
	//
	//// 创建定时器
	//timer := time.NewTimer(next.Sub(time.Now()))
	//// 计数器
	//count = 1
	//log.Info(time.Now(), "  ", "Starting")
	//log.SetReportCaller(true)
	//
	//// 等待定时器触发，执行函数
	//<-timer.C
	//fmt.Println("Time is up,start to send code")

	// 创建 Ticker,之后8小时扫一次
	ticker := cst.GetTicker()

	// 定时器触发时执行的函数
	go func() {
		for {
			select {
			case <-ticker.C:

				if cst.StartTime.IsZero() || cst.EndTime.IsZero() {
					//sendMsg(globalConf.GroupID.AdminGroupID, "请先设置开始时间和结束时间", bot)
					log.Info("请先设置开始时间和结束时间")
					break
				}
				if isWithinTimeRange() && cst.IsSend {
					fmt.Println(time.Now())
					fmt.Println("Time is up,start to send code")
					sendCode(bot, globalConf.GroupID.UserGroupID)
					count++
					log.Info(time.Now(), "第 ", count, " 次运行")
				}

			}
		}
	}()
	// 处理 Ctrl+C 信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号通知
	<-signals
	//
	// 收到信号后关闭 ticker
	ticker.Stop()

}

// getNextExecutionTime 计算下一次执行任务的时间
//func getNextExecutionTime() time.Time {
//	now := time.Now()
//	var next time.Time
//	if now.Hour() >= startTime || now.Hour() < endTime {
//		next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
//	} else {
//		next = time.Date(now.Year(), now.Month(), now.Day(), startTime, 0, 0, 0, now.Location())
//	}
//	return next
//}

// isWithinTimeRange 检查当前时间是否在12:00至22:00之间
func isWithinTimeRange() bool {
	now := time.Now()
	result := now.After(cst.StartTime) && now.Before(cst.EndTime)
	return result
}

// send 发送兑换码
func sendCode(bot *tgbotapi.BotAPI, groupID int64) {
	if groupID == 0 || bot == nil {
		return
	}
	codes, err := GetRandomCode()
	if err != nil {
		sendMsg(groupID, "获取兑换码失败", bot)
		return
	}
	if len(codes) == 0 {
		sendMsg(globalConf.GroupID.AdminGroupID, "兑换码已用完", bot)
		return
	}
	text := "兑换码:\n" + strings.Join(codes, "\n")
	if cst.ActivityText != "" {
		text = cst.ActivityText + "\n" + text
	}
	sendMsg(groupID, text, bot)
	return
}
