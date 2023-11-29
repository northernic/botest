package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func initConfig() {
	files, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
		panic(err)
	}

	//model := make(map[string]Model)

	err = yaml.Unmarshal(files, &globalConf)
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
		panic(err)
	}

	fmt.Println("读取配置成功")
}

// 初始化日志
func initLog() {
	//初始化log
	log = logrus.New()
	file, err := os.OpenFile(LOG, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Error("Failed to open log file: ", err)
		panic(err)
	} else {
		log.SetOutput(file)
	}
}
func initBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(globalConf.BotToken)
	if err != nil {
		log.Error("bot创建出错，错误信息： " + err.Error())
	}
	bot.Debug = true
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)
	userStates = make(map[int64]*UserState)
}
func initRedis() {
	redisCfg := globalConf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error("redis connect ping failed, err:", err)
		fmt.Errorf("redis connect ping failed, err:%v", err)
	} else {
		log.Info("redis connect ping response: pong = ", pong)
		rd = client
	}
}
