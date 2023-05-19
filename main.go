package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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
	Conf = Config{}
	LOG  = "logrus.log"
	log  *logrus.Logger
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
		//wh, _ := tgbotapi.NewWebhook(v)
		//wh.MaxConnections = 10
		//wh.DropPendingUpdates = true
		//apiResponse, err := bot.Request(wh)
		//time.Sleep(time.Second)
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

		//单重命令，示例  /hello
		cmd := update.Message.Command()
		switch cmd {
		case "hello":
			sendMsg(update.Message.Chat.ID, "hello,world!", bot)
		}

		//记录请求
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		arr := strings.Split(update.Message.Text, "/")
		if len(arr) != 0 && arr[0] == "" {
			switch arr[1] {
			//指定显示某模块备用域名
			case "正式域名":
				if len(arr) > 2 && arr[2] != "" {
					//标记是否找到对应模块
					sign := false
					//遍历配置文件，信息匹配
					t := reflect.TypeOf(Conf.Alternate)
					v := reflect.ValueOf(Conf.Alternate)
					for i := 0; i < t.NumField(); i++ {
						value := v.Field(i).Interface()
						s, ok := value.(struct {
							Name          string   `yaml:"name"`
							NewDomainName []string `yaml:"newDomainName"`
							DomainName    []string `yaml:"domainName"`
						})
						if ok {
							if arr[2] == s.Name {
								var text string
								if len(s.NewDomainName) == 0 {
									text = "该正式域名未配置"
								} else {
									text = strings.Join(s.NewDomainName, "  ")
								}
								sendMsg(update.Message.Chat.ID, text, bot)
								sign = true
								break
							}
						} else {
							fmt.Println("请检查配置文件设置")
						}
					}
					if !sign {
						sendMsg(update.Message.Chat.ID, "未找到对应模块，请检查输入或配置文件", bot)
					}
				} else {
					sendMsg(update.Message.Chat.ID, "请输入类型,格式："+"/{备用or正式域名}/{模块名}", bot)
				}
				break
				//指定显示某模块备用域名
			case "备用域名":
				if len(arr) > 2 && arr[2] != "" {
					//标记是否找到对应模块
					sign := false
					//遍历配置文件，信息匹配
					t := reflect.TypeOf(Conf.Alternate)
					v := reflect.ValueOf(Conf.Alternate)
					for i := 0; i < t.NumField(); i++ {
						value := v.Field(i).Interface()
						s, ok := value.(struct {
							Name          string   `yaml:"name"`
							NewDomainName []string `yaml:"newDomainName"`
							DomainName    []string `yaml:"domainName"`
						})
						if ok {
							if arr[2] == s.Name {
								var text string
								if len(s.DomainName) == 0 {
									text = "该备用域名未配置"
								} else {
									text = strings.Join(s.DomainName, "  ")
								}
								sendMsg(update.Message.Chat.ID, text, bot)
								sign = true
								break
							}
						} else {
							fmt.Println("请检查配置文件设置")
						}
					}
					if !sign {
						sendMsg(update.Message.Chat.ID, "未找到对应模块，请检查输入或配置文件", bot)
					}
				} else {
					sendMsg(update.Message.Chat.ID, "请输入类型,格式："+"/备用域名/{模块名}", bot)
				}
				break
			case "groupID":
				sendMsg(update.Message.Chat.ID, "groupID: "+strconv.Itoa(int(update.Message.Chat.ID)), bot)
				break

				//返回上葡京所有域名
			case "上葡京域名":
				text := Conf.ShangPuJing
				sendMsg(update.Message.Chat.ID, text, bot)

				//返回金沙所有域名
			case "金沙域名":
				text := Conf.JinSha
				sendMsg(update.Message.Chat.ID, text, bot)
			case "总列表":
				if len(arr) > 2 && arr[2] != "" {
					switch arr[2] {
					case "aaa":

					}
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
