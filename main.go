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
	"strings"
	"syscall"
	"time"
)

type Config struct {
	DomainName []string `yaml:"domainName"`
	GroupID    int64    `yaml:"groupID"`
	BotToken   string   `yaml:"botToken"`
}

var (
	Conf = Config{}
	LOG  = "logrus.log"
	log  *logrus.Logger
)

func initConfig() {
	files, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
		fmt.Println("读取配置失败")
	}
	err = yaml.Unmarshal(files, &Conf)
	if err != nil {
		log.Fatal(err)
		fmt.Println("读取配置失败")
	}
	fmt.Println("读取配置成功")
}

var count int

func main() {
	initConfig()
	//获取当前时间
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
	//<-timer.C
	fmt.Println("Starting")
	log = logrus.New()
	file, err := os.OpenFile(LOG, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Error("Failed to open log file: ", err)
	} else {
		log.SetOutput(file)
	}
	log.Info(time.Now(), "  ", "Starting")
	//设置后可以在输出日志中显示文件名和方法信息
	log.SetReportCaller(true)

	CheckDomain()
	// 等待定时器触发，执行函数
	<-timer.C
	fmt.Println("Starting")
	CheckDomain()

	// 创建 Ticker
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

	// 收到信号后关闭 ticker
	ticker.Stop()

	// 程序退出
	log.Info(time.Now(), "程序退出 ", "总运行次数: ", count, " 次")
	fmt.Println("Program exiting...")
	os.Exit(0)
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
			tmpMsg = append(tmpMsg, "域名解析出错,该域名为： "+v+"   \n")
			log.Error("域名解析出错,该域名为： " + v)
		} else {
			fmt.Println("域名 " + v + " 信息正常")
		}
	}
	if len(tmpMsg) != 0 {
		_msg := strings.Join(tmpMsg, " ")
		fmt.Println(_msg)
		//msg := tgbotapi.NewMessage(Conf.GroupID, _msg)
		//_, err = bot.Send(msg)
		//if err != nil {
		//	log.Error("bot发送信息出错，错误信息： " + err.Error())
		//}
		fmt.Println("域名解析完毕,记录域名错误成功")
	} else {
		fmt.Println("域名解析完毕,域名信息正常")
		log.Info(time.Now(), " 域名信息正常")
	}
}
