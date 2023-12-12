package main

import (
	"bot/cst"
	"bot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"time"
)

func startTask(bot *tgbotapi.BotAPI, config *TaskConfig) error {
	if bot == nil || config == nil {
		return fmt.Errorf("bot or config is nil")
	}
	if config.EntryID != 0 {
		return nil
	}
	err := checkTaskConfig(config)
	if err != nil {
		return err
	}
	interval := config.Interval
	config.GroupID = globalConf.GroupID.UserGroupID
	entryID, err := Cron.AddFunc("@every "+interval, func() {
		SendCode(config)
	})

	if err != nil {
		log.Error("设置定时任务错误", err)
		sendMsg(globalConf.GroupID.AdminGroupID, "设置定时任务错误", bot)
		return err
	}

	sendMsg(globalConf.GroupID.AdminGroupID, "定时任务已设定,时间间隔为"+interval+"任务id : "+strconv.Itoa(int(entryID)), bot)

	//清空任务配置
	taskConfigs = make(map[int64]*TaskConfig)
	return nil
}

func checkTaskConfig(config *TaskConfig) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if config.Interval == "" {
		return fmt.Errorf("interval is empty")
	}
	if config.CodeNum == 0 {
		return fmt.Errorf("code num is 0")
	}
	if config.StartTime.IsZero() {
		return fmt.Errorf("start time is nil")
	}
	if config.EndTime.IsZero() {
		return fmt.Errorf("end time is nil")
	}
	return nil
}

func sendCodeTimer2(bot *tgbotapi.BotAPI) {

	//interval := cst.SendCodeInterval
	//entryID, err := Cron.AddFunc("@every "+interval, func() {
	//	SendCode(entryID)
	//})
	//if err != nil {
	//	log.Error("设置定时任务错误", err)
	//	sendMsg(globalConf.GroupID.AdminGroupID, "设置定时任务错误", bot)
	//}
	//
	//sendMsg(globalConf.GroupID.AdminGroupID, "定时任务已设定,时间间隔为"+interval+"任务id : "+strconv.Itoa(int(entryID)), bot)

}

func SendCode(config *TaskConfig) {
	if config == nil {
		log.Info(time.Now(), "定时任务配置为空")
		return
	}
	canGo := time.Now().After(config.StartTime) && time.Now().Before(config.EndTime)
	if canGo {
		fmt.Println(time.Now())
		fmt.Println("Time is up,start to send code")
		groupID := config.GroupID
		if groupID == 0 || bot == nil {
			return
		}
		codes, err := GetRandomCode(config.CodeNum)
		if err != nil {
			sendMsg(groupID, "获取兑换码失败", bot)
			return
		}
		if len(codes) == 0 {
			sendMsg(globalConf.GroupID.AdminGroupID, "兑换码已用完", bot)
			return
		}
		text := "兑换码:\n" + strings.Join(codes, "\n")
		if config.ActivityText != "" {
			text = config.ActivityText + "\n" + "\n" + text
		}
		sendMsg(groupID, text, bot)
		return

	}
	return
}

func stopCodeTimer(entryID int) {
	Cron.Remove(cron.EntryID(entryID))
}

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

func GetRandomCode(num int64) ([]string, error) {
	codes, err := rd.SPopN(context.Background(), GiftCodeKey, num).Result()
	if err != nil {
		log.Error("获取随机码失败,err =", err)
		return nil, err
	}
	return codes, nil

}

func SetRandomCode(codeStr string) error {
	codes := utils.StringToStringSli(codeStr)
	for _, v := range codes {
		err := rd.SAdd(context.Background(), GiftCodeKey, v).Err()
		if err != nil {
			log.Error("添加随机码失败,err =", err)
			return err
		}
	}
	return nil
}

func DelRandomCode() error {
	err := rd.Del(context.Background(), GiftCodeKey).Err()
	if err != nil {
		log.Error("删除随机码失败,err =", err)
		return err
	}
	return nil
}
