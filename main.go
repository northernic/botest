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
)

func initConfig() {
	files, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("读取配置失败,err: ", err.Error())
	}
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

func main() {
	initConfig()
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
	bot, err := tgbotapi.NewBotAPI(Conf.BotToken)
	if err != nil {
		log.Error("生成bot接口错误，错误信息： " + err.Error())
		fmt.Println("生成bot接口错误，错误信息： " + err.Error())
	}
	bot.Buffer = 200
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	tmpMsg := []string{}
	for _, v := range Conf.DomainName {
		timeout := 3 * time.Second
		client := http.Client{
			Timeout: timeout,
		}
		fmt.Println("正在访问： ", v)
		resp, err := client.Get(v)
		if err != nil || resp.StatusCode != 200 {
			tmpMsg = append(tmpMsg, "域名解析出错,该域名为： "+v+"\n")
			log.Error("域名解析出错,该域名为： " + v)
			fmt.Println("域名 " + v + " 信息异常")
		} else {
			fmt.Println("域名 " + v + " 信息正常")
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
	bot, err := tgbotapi.NewBotAPI(Conf.BotToken)
	if err != nil {
		log.Error("bot创建出错，错误信息： " + err.Error())
	}
	// 设置机器人接收更新的方式
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	// 处理接收到的更新
	for update := range updates {
		if update.Message == nil { // 忽略非文本消息
			continue
		}
		//仅开头为"/"才处理
		//单重命令(英文)，示例  /hello
		cmd := update.Message.Command()
		if len(cmd) != 0 {
			switch cmd {
			case "hello":
				sendMsg(update.Message.Chat.ID, "hello,world!", bot)
				continue
			case "groupID":
				sendMsg(update.Message.Chat.ID, "groupID: "+strconv.Itoa(int(update.Message.Chat.ID)), bot)
				continue
			case "check":
				CheckDomain()
				continue
			default:
				cmdlist := []string{
					"命令列表大全:",
					"/hello",
					"/check", //检查域名
					"/groupID",
					"/show/{模块名称}",
					"/change/{模块名称}",
					"/正式域名",
					"/备用域名",
					"/上葡京域名",
					"/金沙域名",
					"模块名称:",
					"{ICEX,M1F,LSEX,MIAX,TGX,VGX,ISE,BitBank,SZ,Shop,LuHai}",
				}
				text := strings.Join(cmdlist, "\n")
				sendMsg(update.Message.Chat.ID, text, bot)
				continue
			}
		}

		//记录请求
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//多重命令  示例 /show/iex
		arr := strings.Split(update.Message.Text, "/")
		if len(arr) != 0 && arr[0] == "" {
			switch arr[1] {
			case "上葡京域名":
				text := Conf.ShangPuJing
				sendMsg(update.Message.Chat.ID, text, bot)
				//返回金沙所有域名
			case "金沙域名":
				text := Conf.JinSha
				sendMsg(update.Message.Chat.ID, text, bot)
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
						if field.Name == arr[2] {
							fieldValue := v.FieldByName(field.Name)
							sign = true
							text := getFieldInfo(fieldValue)
							sendMsg(update.Message.Chat.ID, text, bot)
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
			}
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
