package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Conf       = Config{}
	LOG        = "logrus.log"
	log        *logrus.Logger
	configName = "config.yaml"
	bot        *tgbotapi.BotAPI
	userStates map[int64]*UserState
)

type UserState struct {
	Uid               int
	LastCallbackMsgID int
	LastCallbackData  string
	ErrorCode         string
	//Sign              bool //true代表已处理
}

func initConfig() {
	files, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
	}

	//model := make(map[string]Model)

	err = yaml.Unmarshal(files, &Conf)
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
	}
	//初始化log
	log = logrus.New()
	file, err := os.OpenFile(LOG, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Error("Failed to open log file: ", err)
	} else {
		log.SetOutput(file)
	}
	fmt.Println("读取配置成功")
}

var count int

func initBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(Conf.BotToken)
	if err != nil {
		log.Error("bot创建出错，错误信息： " + err.Error())
	}
	bot.Debug = true
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)
	userStates = make(map[int64]*UserState)
}

func main() {
	initConfig()
	initBot()
	check := true //开关域名扫描
	go startBot()
	if check {
		//启动先扫描一遍
		CheckDomain()

		//八小时跑一次
		now := time.Now()
		var next time.Time
		if now.Hour() < 24 {
			next = time.Date(now.Year(), now.Month(), now.Day(), (now.Hour()/8+1)*8, 0, 0, 0, now.Location())
		} else {
			next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		}
		// 创建定时器
		timer := time.NewTimer(next.Sub(now))
		// 计数器
		count = 1
		log.Info(time.Now(), "  ", "Starting")
		log.SetReportCaller(true)
		// 等待定时器触发，执行函数
		<-timer.C
		fmt.Println("Time is up,start Checking")
		CheckDomain()

		// 创建 Ticker,之后8小时扫一次
		ticker := time.NewTicker(8 * time.Hour)

		// 定时器触发时执行的函数
		go func() {
			for {
				select {
				case <-ticker.C:
					fmt.Println(time.Now())
					CheckDomain()
					count++
					log.Info(time.Now(), "第 ", count, " 次运行")
				}
			}
		}()
		// 处理 Ctrl+C 信号
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		// 等待信号通知
		<-signals
		//
		// 收到信号后关闭 ticker
		ticker.Stop()
	}
	select {}
}

func CheckDomain() {
	tmpMsg := []string{}
	if len(Conf.DomainName) == 0 {
		return
	}
	for _, v := range Conf.DomainName {
		timeout := 3 * time.Second
		client := http.Client{
			Timeout: timeout,
		}
		fmt.Println("正在访问： ", v)
		resp, err := client.Get(v)
		if err != nil {
			tmpMsg = append(tmpMsg, "访问出错，该域名为："+v+"\n")
			log.Error("访问出错，该域名为：" + v)
			fmt.Println("域名 " + v + " 信息异常")
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				tmpMsg = append(tmpMsg, "状态码异常，该域名为："+v+"\n")
				log.Error("状态码异常，该域名为：" + v)
				fmt.Println("域名 " + v + " 信息异常")
			} else {
				fmt.Println("域名 " + v + " 信息正常")
			}
		}
	}
	l := len(tmpMsg)
	if l != 0 {
		//10条错误发送一次tel
		if l <= 10 {
			sendMsg(Conf.GroupID, strings.Join(tmpMsg, " "), bot)
		} else {
			for i := 0; i < l; i += 10 {
				end := i + 10
				if end > l {
					end = l
				}
				sendMsg(Conf.GroupID, strings.Join(tmpMsg[i:end], " "), bot)
			}
		}
		fmt.Println("域名解析完毕,记录域名错误成功")
	} else {
		fmt.Println("域名解析完毕,域名信息正常")
		log.Info(time.Now(), " 域名信息正常")
	}
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
	//updates, err := bot.GetUpdates(u)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if len(updates) > 0 {
	//	lastUpdate := updates[len(updates)-1]
	//	offset := lastUpdate.UpdateID + 1
	//	// 使用 offset 设置下一次获取的起始位置
	//	u.Offset = offset
	//}
	u.Timeout = 60
	updateChan, _ := bot.GetUpdatesChan(u)

	// 处理接收到的更新
	for update := range updateChan {
		if update.Message == nil { // 忽略非文本消息
			continue
		}
		if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
			continue
		}
		//仅开头为"/"才处理
		//单重命令(英文)，示例  /hello
		cmd := update.Message.Command()

		cmd = strings.ToLower(cmd)
		if len(cmd) != 0 {
			switch cmd {
			case "hello":
				sendMsg(update.Message.Chat.ID, "hello,world!", bot)
				continue
			case "groupid":
				sendMsg(update.Message.Chat.ID, "groupID: "+strconv.Itoa(int(update.Message.Chat.ID)), bot)
				continue
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
				continue
			case "change":
				nextLevelInlineKeyboard := packDomainKeyboard(Conf.Domains)
				reply := "选择要修改域名的模块："
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = nextLevelInlineKeyboard
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				continue
			case "check":
				CheckDomain()
				continue
			case "list":
				cmdlist := []string{
					"命令列表大全:",
					"/hello",
					"/check", //检查域名
					"/groupID",
					"/myID",
					"/show/{模块名称}",
					"/change/{模块名称}",
					"/add/",
					"/delete/",
					"/remove",
					"/上葡京域名",
					"/金沙域名",
					"模块名称：{ICEX,M1F,MIAX,TGX,VGX,ISE,BitBank,Shop,Voya}",
				}
				strconv.FormatInt(4, 2)
				text := strings.Join(cmdlist, "\n")
				sendMsg(update.Message.Chat.ID, text, bot)
				continue
			default:
				continue
			}

		}

		//记录请求
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//多重命令  示例 /show/iex
		arr := strings.Split(update.Message.Text, "/")
		if len(arr) != 0 && arr[0] == "" {
			if len(arr) == 1 {
				continue
			}
			switch arr[1] {
			case "上葡京域名":
				text := Conf.ShangPuJing
				result := checkAuth(update.Message.Chat.ID, "shangpujing")
				if result {
					sendMsg(update.Message.Chat.ID, text, bot)
					continue
				}
				sendMsg(update.Message.Chat.ID, "权限不足", bot)

			case "金沙域名":
				text := Conf.JinSha
				result := checkAuth(update.Message.Chat.ID, "jinsha")
				if result {
					sendMsg(update.Message.Chat.ID, text, bot)
					//返回金沙所有域名
					continue
				}
				sendMsg(update.Message.Chat.ID, "权限不足", bot)
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
			case "add":
				if len(arr) > 2 && arr[2] != "" {
					switch arr[2] {

					}
				}
			case "delete":
				if len(arr) > 2 && arr[2] != "" {

				}

			case "setdomain":
				//向指定url发送https post 请求
				//bot.SetWebhook()

				//测试修改config文件
			//case "test":
			//	Conf.JinSha = "testModify"
			//	configData, err := yaml.Marshal(&Conf)
			//	if err != nil {
			//		fmt.Printf("Error marshaling config data: %s\n", err)
			//		continue
			//	}
			//	err = ioutil.WriteFile("config.yaml", configData, 0644)
			//	if err != nil {
			//		fmt.Printf("Error marshaling config data: %s\n", err)
			//		continue
			//	}
			//	sendMsg(update.Message.Chat.ID, "修改成功", bot)
			default:
				sendMsg(update.Message.Chat.ID, "请输入类型,格式："+"/{命令}/{模块名}", bot)
				continue
			}
		}

		// 处理用户的文本输入，可以根据需要进行逻辑处理
		//reply := "收到您的输入：" + update.Message.Text
		//
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		//_, err := bot.Send(msg)
		//if err != nil {
		//	log.Println(err)
		//}
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
	for _, v := range Conf.GroupAuth {
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
	case "luhai":
		return LuHai
	case "voya":
		return Voya
	case "jinsha":
		return JinSha
	case "shangpujing":
		return ShangPuJing
	default:
		return 0
	}

}
