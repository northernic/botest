package main

import (
	"bot/cst"
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
	case "hello":
		sendMsg(update.Message.Chat.ID, "hello,world!", bot)
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
		sendMsg(update.Message.Chat.ID, "兑换码开始发放", bot)
		break
	case "stopsendcode":
		cst.IsSend = false
		sendMsg(update.Message.Chat.ID, "兑换码停止发放", bot)
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
	case "list":
		text := strings.Join(cst.CmdList, "\n")
		sendMsg(update.Message.Chat.ID, text, bot)
		break
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
		sendMsg(update.Message.Chat.ID, "设置结束时间成功", bot)
		break
	case "setactivitytext":
		text := arr[2]
		cst.ActivityText = text
		sendMsg(update.Message.Chat.ID, "设置活动文案成功", bot)
		break
	default:
		//sendMsg(update.Message.Chat.ID, "请输入类型,格式："+"/{命令}/{模块名}", bot)
		break
	}
}
