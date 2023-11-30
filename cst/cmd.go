package cst

var CmdList = []string{
	"命令列表大全:",
	"/hello ",
	"/check", //检查域名
	"/groupID 显示当前群id",
	"/startSendCode 开始发放兑换码",
	"/stopSendCode 停止发放兑换码",
	"/setstartTime/2006-01-02 15:04:05 设置兑换码发放开始时间",
	"/setendTime/2006-01-02 15:04:05 设置兑换码发放结束时间",
	"/setactivitytext/文案 设置活动文案",
	"/getcode 获取五个兑换码",
	"/setcode/{兑换码1 兑换码2 兑换码3 ...}    兑换码之间用换行隔开",
	"/delcode 清空兑换码",
	"/myID 显示当前用户id",
	"/show/{模块名称}",
	"/change/{模块名称}",
	"/add/",
	"/delete/",
	"/remove",
	"模块名称：{}",
}
