package main

import (
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	globalConf = Config{}
	LOG        = "logrus.log"
	log        *logrus.Logger
	configName = "config.yaml"
	bot        *tgbotapi.BotAPI
	rd         *redis.Client
	userStates map[int64]*UserState

	Cron        *cron.Cron
	CronEntries = make(map[cron.EntryID]*cron.Cron)

	TaskList = make(map[int]*TaskConfig)
)

type UserState struct {
	Uid               int
	LastCallbackMsgID int
	LastCallbackData  string
	ErrorCode         string
	//Sign              bool //true代表已处理
}

var (
	count       int
	GiftCodeKey = "GiftCode"
)

type TaskConfig struct {
	EntryID        cron.EntryID `json:"entry_id,omitempty"`        //任务id
	Interval       string       `json:"interval,omitempty"`        //间隔时间  1h2m3s
	CodeNum        int64        `json:"code_num,omitempty"`        //兑换码数量
	IsSend         bool         `json:"is_send,omitempty"`         //是否发送
	ActivityText   string       `json:"activity_text,omitempty"`   //活动文案
	StartTime      time.Time    `json:"start_time"`                // 每天的开始时间（小时）
	EndTime        time.Time    `json:"end_time"`                  // 次日的结束时间（小时）
	GroupID        int64        `json:"group_id,omitempty"`        //群组id
	CurrentSetting string       `json:"current_setting,omitempty"` //当前设置参数字段
}

var taskConfigs = make(map[int64]*TaskConfig)
