package main

func CheckDomain() {
	//tmpMsg := []string{}
	//if len(globalConf.DomainName) == 0 {
	//	return
	//}
	//for _, v := range globalConf.DomainName {
	//	timeout := 3 * time.Second
	//	client := http.Client{
	//		Timeout: timeout,
	//	}
	//	fmt.Println("正在访问： ", v)
	//	resp, err := client.Get(v)
	//	if err != nil {
	//		tmpMsg = append(tmpMsg, "访问出错，该域名为："+v+"\n")
	//		log.Error("访问出错，该域名为：" + v)
	//		fmt.Println("域名 " + v + " 信息异常")
	//	} else {
	//		defer resp.Body.Close()
	//		if resp.StatusCode != http.StatusOK {
	//			tmpMsg = append(tmpMsg, "状态码异常，该域名为："+v+"\n")
	//			log.Error("状态码异常，该域名为：" + v)
	//			fmt.Println("域名 " + v + " 信息异常")
	//		} else {
	//			fmt.Println("域名 " + v + " 信息正常")
	//		}
	//	}
	//}
	//l := len(tmpMsg)
	//if l != 0 {
	//	//10条错误发送一次tel
	//	if l <= 10 {
	//		sendMsg(globalConf.GroupID, strings.Join(tmpMsg, " "), bot)
	//	} else {
	//		for i := 0; i < l; i += 10 {
	//			end := i + 10
	//			if end > l {
	//				end = l
	//			}
	//			sendMsg(globalConf.GroupID, strings.Join(tmpMsg[i:end], " "), bot)
	//		}
	//	}
	//	fmt.Println("域名解析完毕,记录域名错误成功")
	//} else {
	//	fmt.Println("域名解析完毕,域名信息正常")
	//	log.Info(time.Now(), " 域名信息正常")
	//}
}
