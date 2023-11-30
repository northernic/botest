package main

import (
	"bot/cst"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"time"
)

func sendCodeTimer2(bot *tgbotapi.BotAPI) {

	interval := cst.SendCodeInterval
	entryID, err := Cron.AddFunc("@every "+interval, Send)
	if err != nil {
		log.Error("设置定时任务错误", err)
		sendMsg(globalConf.GroupID.AdminGroupID, "设置定时任务错误", bot)
	}

	sendMsg(globalConf.GroupID.AdminGroupID, "定时任务已设定,时间间隔为"+interval+"任务id : "+strconv.Itoa(int(entryID)), bot)

}

func Send() {
	if isWithinTimeRange() && cst.IsSend {
		fmt.Println(time.Now())
		fmt.Println("Time is up,start to send code")
		sendCode(bot, globalConf.GroupID.UserGroupID)
		count++
		log.Info(time.Now(), "第 ", count, " 次运行")
	}
	return
}

func stopCodeTimer(entryID int) {
	Cron.Remove(cron.EntryID(entryID))
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
	codes, err := GetRandomCode(cst.CodeNum)
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
