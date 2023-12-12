package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// 初始化
func init() {
	initConfig()     //初始化配置文件
	initLog()        //初始化日志
	initBot()        //初始化bot
	initRedis()      //初始化redis
	initCron()       //初始化定时任务
	initTaskParams() //初始化任务参数
}

func main() {

	//check := true //开关域名扫描
	go startBot()

	select {}
}

// 发送消息给指定聊天ID
func sendMsg(chatID int64, msg string, bot *tgbotapi.BotAPI) {
	if msg == "" {
		return
	}
	tgMsg := tgbotapi.NewMessage(chatID, msg)
	_, err := bot.Send(tgMsg)
	if err != nil {
		log.Error("bot发送信息出错，错误信息： " + err.Error())
	}
}

func startBot() {
	//设置机器人接收更新的方式
	u := tgbotapi.NewUpdate(0)
	//这里注释的是只处理最新的更新
	updates, err := bot.GetUpdates(u)
	if err != nil {
		log.Fatal(err)
	}
	if len(updates) > 0 {
		lastUpdate := updates[len(updates)-1]
		offset := lastUpdate.UpdateID + 1
		// 使用 offset 设置下一次获取的起始位置
		u.Offset = offset
	}
	u.Timeout = 60
	updateChan, _ := bot.GetUpdatesChan(u)

	// 处理接收到的更新
	for update := range updateChan {
		//if update.Message == nil { // 忽略非文本消息
		//	continue
		//}
		if update.CallbackQuery != nil {
			//处理回调查询
			handleCallback(update.CallbackQuery)
			continue
		}

		// 处理非命令的文本消息（用于设置时间、间隔等）
		if update.Message.Text != "" && !update.Message.IsCommand() {
			chatID := update.Message.Chat.ID
			config, exists := taskConfigs[chatID]
			if exists {
				handleTextResponse(bot, update.Message, config)
			}
		}

		//判断请求来源群
		switch update.Message.Chat.ID {
		case globalConf.GroupID.UserGroupID:
			//用户组
			continue
		case globalConf.GroupID.AdminGroupID:
			//管理员组
			HandleAdminText(update, bot)
			continue
		default:
			log.Info("未知群组请求,groupID = ", update.Message.Chat.ID)
			continue
		}

	}

}

// 调用其他网址的API函数
func callAPI(args string) string {
	// 根据参数访问其他网址的API

	// 执行API请求的逻辑...

	// 返回API响应
	return "API响应"
}
