package main

import (
	"bot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// 初始化
func init() {
	initConfig() //初始化配置文件
	initLog()    //初始化日志
	initBot()    //初始化bot
	initRedis()  //初始化redis
	initCron()   //初始化定时任务
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
		if update.Message == nil { // 忽略非文本消息
			continue
		}
		if update.CallbackQuery != nil {
			//处理回调查询
			handleCallback(update.CallbackQuery)
			continue
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

// 读取 config.yaml 文件并返回 Config 结构体
func readConfigFile() (*Config, error) {
	config := &Config{}

	content, err := ioutil.ReadFile(configName)
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getFieldInfo(value reflect.Value) string {
	typeName := value.Type()
	var st []string
	for i := 0; i < value.NumField(); i++ {
		typeField := typeName.Field(i)
		fieldName := typeField.Name
		fieldValue := value.Field(i).Interface()

		// 处理切片类型
		if value.Field(i).Kind() == reflect.Slice {
			sliceValues := make([]string, value.Field(i).Len())
			for j := 0; j < value.Field(i).Len(); j++ {
				sliceValues[j] = fmt.Sprintf("%v", value.Field(i).Index(j))
			}
			fieldValue = strings.Join(sliceValues, "\n")
		}
		//仅展示配置项
		if fieldValue != "" {
			tmpSt := fmt.Sprintf("%s:\n%v\n", fieldName, fieldValue)
			st = append(st, tmpSt)
		}
	}
	return strings.Join(st, "\n")
}

// 调用其他网址的API函数
func callAPI(args string) string {
	// 根据参数访问其他网址的API

	// 执行API请求的逻辑...

	// 返回API响应
	return "API响应"
}

func packDomainKeyboard(domains map[string]int64) tgbotapi.InlineKeyboardMarkup {
	nextLevelInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	row := []tgbotapi.InlineKeyboardButton{}
	for k := range domains {
		button := tgbotapi.NewInlineKeyboardButtonData(k, k)
		row = append(row, button)
		if len(row) == 3 {
			nextLevelInlineKeyboard.InlineKeyboard = append(nextLevelInlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	// 如果最后一行只有一个按钮，将其添加到内联键盘
	if len(row) == 1 {
		nextLevelInlineKeyboard.InlineKeyboard = append(nextLevelInlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(row...))
	}
	return nextLevelInlineKeyboard
}

func handleCallback(callback *tgbotapi.CallbackQuery) {

	switch callback.Data {
	case "盘口":
		// 生成选项一的下一层内联键盘
		nextLevelInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("601", "601"), // 错误码后续在这里更新，并增加case的处理
			),
		)
		// 更新原始消息的内联键盘为下一层内联键盘
		//editMsg := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, nextLevelInlineKeyboard)
		//_, err := bot.Send(editMsg)

		reply := "选择错误码："
		editMsgText := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, reply)
		editMsgText.ReplyMarkup = &nextLevelInlineKeyboard // 设置新的内联键盘
		_, err := bot.Send(editMsgText)
		if err != nil {
			log.Println(err)
		}

	default:
		// 处理未知的回调查询数据
	}
}

func checkAuth(groupID int64, moduleName string) bool {
	//先查找groupID权限
	authID := ""
	for _, v := range globalConf.GroupAuth {
		if v.ID == groupID {
			authID = v.AuthID
		}
	}
	//查找模块需要什么权限
	moduleAuthID := getmoduleAuthID(moduleName)

	if authID != "" && moduleAuthID != 0 {
		authid, err := strconv.ParseInt(authID, 2, 64)
		if err != nil {
			return false
		}
		return int(authid)&moduleAuthID == moduleAuthID
	}
	return false
}

func getmoduleAuthID(moduleName string) int {
	moduleName = strings.ToLower(moduleName)
	switch moduleName {
	case "icex":
		return ICEX
	case "m1f":
		return M1F
	case "miax":
		return MIAX
	case "tgx":
		return TGX
	case "vgx":
		return VGX
	case "ise":
		return ISE
	case "bitbank":
		return BitBank
	case "sz":
		return SZ
	case "shop":
		return Shop
	case "aquis":
		return Aquis
	case "voya":
		return Voya
	case "jinsha":
		return JinSha
	case "shangpujing":
		return ShangPuJing
	case "jason":
		return Jason
	default:
		return 0
	}

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

// mentionBot检查消息是否包含对机器人的提及
func mentionBot(message *tgbotapi.Message, botUsername string) bool {
	if message == nil {
		return false
	}
	if message.Entities == nil {
		return false
	}
	for _, entity := range *message.Entities {
		if entity.Type == "mention" {
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			if mention == "@"+botUsername {
				return true
			}
		}
	}
	return false
}
