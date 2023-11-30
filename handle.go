package main

import (
	"bot/cst"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
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
			cmdlist := []string{
				"测试命令列表大全:",
				"/hello",
				//"/check", //检查域名
				"/groupID",
				"/getcode",
				"/setcode/{随机码1;随机码2;随机码3;...}",
				"/myID",
				//"/show/{模块名称}",
				//"/change/{模块名称}",
				//"/add/",
				//"/delete/",
				//"/remove",
				//"模块名称：{}",
			}
			strconv.FormatInt(4, 2)
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
