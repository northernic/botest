package main

type DashboardType string

const (
	DashboardTypeStartTime    string = "startTime"        // 开始时间
	DashboardTypeEndTime      string = "endTime"          // 结束时间
	DashboardTypeInterval     string = "sendCodeInterval" // 间隔时间
	DashboardTypeCodeNum      string = "codeNum"          // 兑换码数量
	DashboardTypeIsSend       string = "isSend"           // 是否发送
	DashboardTypeActivityText string = "activityText"     // 活动文案

	DashboardTypeStartTask string = "start_task" // 开始任务

)
