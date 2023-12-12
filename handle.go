package main

import (
	"bot/cst"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func HandleAdminText(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	//记录请求
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	cmdType := cst.GetCmdType(update.Message.Text)
	switch cmdType {
	case cst.CmdTypeSingle:
		//处理单独命令 单重命令(英文)，示例  /hello
		HandleCmd(update, strings.ToLower(update.Message.Command()))
		break
	case cst.CmdTypeMul:
		//处理多重命令  示例 /show/iex
		HandleMulCmd(update, strings.Split(update.Message.Text, "/"))
		break
	default:
		break
	}

	if update.Message.Chat.Type == cst.ChatTypeGroup {
		if mentionBot(update.Message, bot.Self.UserName) {
			replyText := "你在叫我吗？\n"
			cmdlist := cst.CmdList
			text := strings.Join(cmdlist, "\n")
			sendMsg(update.Message.Chat.ID, replyText+text, bot)
			return
		}

	}

	//处理用户的文本输入，可以根据需要进行逻辑处理
	//reply := "收到您的输入：" + update.Message.Text
	//sendMsg(update.Message.Chat.ID, reply, bot)
	return
}

// 处理单独命令
func HandleCmd(update tgbotapi.Update, cmd string) {
	switch cmd {
	case "list":
		text := strings.Join(cst.CmdList, "\n")
		sendMsg(update.Message.Chat.ID, text, bot)
		break
	case "hello":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "请选择模块"+"\n")
		msg.ReplyMarkup = mainKeyboard()
		bot.Send(msg)
		break
	case "groupid":
		sendMsg(update.Message.Chat.ID, "groupID: "+strconv.Itoa(int(update.Message.Chat.ID)), bot)
		break
	case "remove":
		// 创建一个发送给用户的空的ReplyKeyboardRemove
		removeKeyboard := tgbotapi.NewRemoveKeyboard(false)

		// 设置要移除键盘的目标聊天ID
		chatID := int64(update.Message.Chat.ID) // 替换为实际的聊天ID,id为群id的时候是清除全部群成员的自定义键盘

		// 替换为实际的聊天ID,个人的话可以用这个
		//chatID := int64(update.Message.From.ID)

		// 创建一个新的消息配置
		msg := tgbotapi.NewMessage(chatID, "移除自定义键盘")
		msg.ReplyMarkup = removeKeyboard

		// 发送消息
		_, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	case "myid":
		sendMsg(update.Message.Chat.ID, "myID: "+strconv.Itoa(update.Message.From.ID), bot)
		break
	case "change":
		nextLevelInlineKeyboard := packDomainKeyboard(globalConf.Domains)
		reply := "选择要修改域名的模块："
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyMarkup = nextLevelInlineKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		break
	case "check":
		CheckDomain()
		break
	case "startsendcode":
		if cst.StartTime.IsZero() || cst.EndTime.IsZero() {
			sendMsg(update.Message.Chat.ID, "请先设置开始时间和结束时间", bot)
			break
		}
		cst.IsSend = true
		sendCodeTimer2(bot)
		sendMsg(update.Message.Chat.ID, "兑换码开始发放", bot)
		break
	case "getcode":
		sendCode(bot, update.Message.Chat.ID)
		break
	case "delcode":
		err := DelRandomCode()
		if err != nil {
			sendMsg(update.Message.Chat.ID, "清空兑换码失败", bot)
			break
		}
		sendMsg(update.Message.Chat.ID, "清空兑换码成功", bot)
		break
	case "settask":
		//设置任务
		chatID := update.Message.Chat.ID
		taskConfigs[chatID] = &TaskConfig{}
		msg := tgbotapi.NewMessage(chatID, "Let's set up your task. What would you like to set first?")
		msg.ReplyMarkup = taskSettingKeyboard()

		bot.Send(msg)

	default:
		break
	}
}

func HandleMulCmd(update tgbotapi.Update, arr []string) {
	if len(arr) == 1 || len(arr) == 2 {
		return
	}
	if arr[1] == "" || arr[2] == "" {
		return
	}
	switch strings.ToLower(arr[1]) {
	case "show":
		if len(arr) > 2 && arr[2] != "" {
			//标记是否找到对应模块
			sign := false
			conf, err := readConfigFile()
			if err != nil {
				log.Error("Failed to open log file: ", err)
				sendMsg(update.Message.Chat.ID, "请检查配置文件", bot)
			}
			//遍历配置文件，信息匹配
			t := reflect.TypeOf(*conf)
			v := reflect.ValueOf(*conf)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if strings.ToLower(field.Name) == strings.ToLower(arr[2]) {
					fieldValue := v.FieldByName(field.Name)
					sign = true
					result := checkAuth(update.Message.Chat.ID, strings.ToLower(arr[2]))
					if result {
						text := getFieldInfo(fieldValue)
						sendMsg(update.Message.Chat.ID, text, bot)
						break
					}
					sendMsg(update.Message.Chat.ID, "权限不足", bot)
					break

				}
			}
			if !sign {
				sendMsg(update.Message.Chat.ID, "未找到对应模块，请检查输入或配置文件", bot)
			}
		}

	case "setcode":
		if len(arr) > 2 && arr[2] != "" {
			err := SetRandomCode(arr[2])
			if err != nil {
				sendMsg(update.Message.Chat.ID, "设置随机码失败", bot)
				break
			}
			sendMsg(update.Message.Chat.ID, "设置兑换码成功", bot)
			break
		}

	case "setstarttime":
		startTimeStr := arr[2]
		//默认输入巴西时间
		startTime, err := time.Parse(time.DateTime, startTimeStr)
		if err != nil {
			log.Error("开始时间格式错误", err)
			sendMsg(update.Message.Chat.ID, "时间格式错误,请按照如下格式\n"+time.DateTime, bot)
			break
		}
		cst.StartTime = startTime.Add(cst.BrazilOffset)
		err = rd.Set(context.Background(), DashboardTypeStartTime, startTimeStr, -1).Err()
		if err != nil {
			log.Error("开始时间保存redis错误", err)
			break
		}
		sendMsg(update.Message.Chat.ID, "设置开始时间成功", bot)
		break
	case "setendtime":
		endTimeStr := arr[2]
		endTime, err := time.Parse(time.DateTime, endTimeStr)
		if err != nil {
			log.Error("结束时间格式错误", err)
			sendMsg(update.Message.Chat.ID, "时间格式错误,请按照如下格式\n"+time.DateTime, bot)
			break
		}
		cst.EndTime = endTime.Add(cst.BrazilOffset)
		err = rd.Set(context.Background(), DashboardTypeEndTime, endTimeStr, -1).Err()
		if err != nil {
			log.Error("结束时间时间保存redis错误", err)
			break
		}
		sendMsg(update.Message.Chat.ID, "设置结束时间成功", bot)
		break
	case "setactivitytext":
		text := arr[2]
		cst.ActivityText = text
		err := rd.Set(context.Background(), DashboardTypeActivityText, text, -1).Err()
		if err != nil {
			log.Error("活动文案保存redis错误", err)
			break
		}
		sendMsg(update.Message.Chat.ID, "设置活动文案成功", bot)
		break
	case "setinterval":
		intervalTime := arr[2]
		if intervalTime == "" {
			sendMsg(update.Message.Chat.ID, "请输入时间间隔", bot)
			break
		}
		isValid := cst.CheckInterval(intervalTime)
		if !isValid {
			sendMsg(update.Message.Chat.ID, "请输入正确的时间间隔(1h2m3s)-1小时2分钟3秒", bot)
			break
		}
		cst.SendCodeInterval = intervalTime
		err := rd.Set(context.Background(), DashboardTypeInterval, intervalTime, -1).Err()
		if err != nil {
			log.Error("时间间隔保存redis错误", err)
			break
		}
		sendMsg(update.Message.Chat.ID, "设置时间间隔成功", bot)
		break
	case "setcodenum":
		numStr := arr[2]
		if numStr == "" {
			sendMsg(update.Message.Chat.ID, "请输入兑换码数量", bot)
			break
		}
		num, err := strconv.Atoi(numStr)
		if err != nil {
			sendMsg(update.Message.Chat.ID, "请输入正确的兑换码数量(纯数字)", bot)
			break
		}
		err = rd.Set(context.Background(), DashboardTypeCodeNum, numStr, -1).Err()
		if err != nil {
			log.Error("兑换码数量保存redis错误", err)
			break
		}
		cst.CodeNum = int64(num)
		sendMsg(update.Message.Chat.ID, "设置兑换码数量成功", bot)
		break
	case "stopsendcode":
		entryIDStr := arr[2]
		if entryIDStr == "" {
			sendMsg(update.Message.Chat.ID, "请输入要停止的任务id", bot)
			break
		}
		entryID, err := strconv.Atoi(entryIDStr)
		if err != nil {
			sendMsg(update.Message.Chat.ID, "请输入要停止的任务id(纯数字)", bot)
			break
		}
		stopCodeTimer(entryID)

		sendMsg(update.Message.Chat.ID, "定时任务 ID:"+entryIDStr+" 已停止", bot)
		break
	default:
		//sendMsg(update.Message.Chat.ID, "请输入类型,格式："+"/{命令}/{模块名}", bot)
		break
	}
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, callback.Data))

	switch callback.Data {
	case DashboardTypeStartTime, DashboardTypeEndTime, DashboardTypeInterval,
		DashboardTypeCodeNum, DashboardTypeActivityText:
		handleDashboardInput(bot, chatID, callback.Data)

	case DashboardTypeStartTask:
		config, exists := taskConfigs[chatID]
		if !exists {
			return
		}
		err := startTask(bot, config)
		if err != nil {
			sendMsg(chatID, "任务启动失败,err = "+err.Error(), bot)
			break
		}

	case cst.KeyboardTypeSettingTask:
		//设置任务
		taskConfigs[chatID] = &TaskConfig{}
		msg := tgbotapi.NewMessage(chatID, "开始设置定时任务,请按照流程设置\n")
		msg.ReplyMarkup = taskSettingKeyboard()
		bot.Send(msg)

	case cst.KeyboardTypeTaskList:
		//任务列表(在cron中查看)
		a := Cron.Entries()
		if len(a) == 0 {
			sendMsg(chatID, "暂无定时任务", bot)
			break
		}
		st := formatCronEntries(a)
		msg := tgbotapi.NewMessage(chatID, st)
		bot.Send(msg)

	case cst.KeyboardTypeMain:
		bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatID, callback.Message.MessageID, mainKeyboard()))
		//msg := tgbotapi.NewMessage(chatID, "请选择模块"+"\n")
		//msg.ReplyMarkup = mainKeyboard()
		//bot.Send(msg)
		break

	case cst.KeyboardTypeCancel:
		//删除这个回调查询
		bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		break

	default:
		// 处理未知的回调查询数据
	}

}

func handleDashboardInput(bot *tgbotapi.BotAPI, chatID int64, dashboardType string) {
	config, exists := taskConfigs[chatID]
	if !exists {
		//处理不存在 config 的情况
		return
	}

	// 使用映射来获取对应的提示消息
	prompt, ok := settingPrompts[dashboardType]
	if !ok {
		// 处理未知 dashboardType 的情况
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, prompt))
	config.CurrentSetting = dashboardType
}

// 处理用户的文本输入，可以根据需要进行逻辑处理
func handleTextResponse(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *TaskConfig) {
	var reply string
	switch config.CurrentSetting {
	case DashboardTypeStartTime:
		// 解析开始时间
		startTime, err := time.Parse(time.DateTime, message.Text)

		if err != nil {
			reply = "无效的时间格式,请按照如下格式重新输入\n" + time.DateTime
		} else {
			config.StartTime = startTime.Add(cst.BrazilOffset)
			config.CurrentSetting = ""
			reply = fmt.Sprintf("开始时间设置成功: %s", startTime.Format(time.DateTime))
		}
	case DashboardTypeEndTime:
		// 解析结束时间
		endTime, err := time.Parse(time.DateTime, message.Text)
		if err != nil {
			reply = "无效的时间格式,请按照如下格式重新输入\n" + time.DateTime
		} else {
			config.EndTime = endTime.Add(cst.BrazilOffset)
			config.CurrentSetting = ""
			reply = fmt.Sprintf("结束时间已设置: %s", endTime.Format(time.DateTime))
		}
	case DashboardTypeInterval:
		// 验证并设置间隔
		if isValidInterval(message.Text) {
			config.Interval = message.Text
			config.CurrentSetting = ""
			reply = fmt.Sprintf("时间间隔已设置: %s", message.Text)
		} else {
			reply = "无效的时间间隔,请按照如下格式重新输入\n" + "1h2m3s"
		}
	case DashboardTypeCodeNum:
		num, err := strconv.Atoi(message.Text)
		if err != nil {
			log.Error("单次发码数量解析失败, err:", err)
			reply = "请输入正确的兑换码数量(纯数字)"
			break
		}
		config.CodeNum = int64(num)
		config.CurrentSetting = ""
		reply = fmt.Sprintf("单次发码数量已设置: %s", message.Text)
	case DashboardTypeActivityText:
		config.ActivityText = message.Text
		config.CurrentSetting = ""
		reply = fmt.Sprintf("活动文案已设置: %s", message.Text)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	if config.CurrentSetting == "" {
		msg.ReplyMarkup = taskSettingKeyboard()
	}
	bot.Send(msg)

}
