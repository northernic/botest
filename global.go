package main

import (
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var (
	globalConf = Config{}
	LOG        = "logrus.log"
	log        *logrus.Logger
	configName = "config.yaml"
	bot        *tgbotapi.BotAPI
	rd         *redis.Client
	userStates map[int64]*UserState
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
