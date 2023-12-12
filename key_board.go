package main

import (
	"bot/cst"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var settingPrompts = map[string]string{
	DashboardTypeStartTime:    "请输入开始时间 (2006-01-02 15:04:05):",
	DashboardTypeEndTime:      "请输入结束时间 (2006-01-02 15:04:05):",
	DashboardTypeInterval:     "请输入时间间隔(格式必须符合1h2m3s):",
	DashboardTypeCodeNum:      "请输入单次发码数量(纯数字):",
	DashboardTypeActivityText: "请输入活动文案:(可以为空)",
}

// 主功能菜单
func mainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("主菜单", cst.KeyboardTypeMain),
			tgbotapi.NewInlineKeyboardButtonData("定时任务列表", cst.KeyboardTypeTaskList),
			tgbotapi.NewInlineKeyboardButtonData("设置一个定时任务", cst.KeyboardTypeSettingTask),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("取消", cst.KeyboardTypeCancel),
		),
	)
}

// 定时任务设定菜单
func taskSettingKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("主菜单", cst.KeyboardTypeMain),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("设置开始时间", DashboardTypeStartTime),
			tgbotapi.NewInlineKeyboardButtonData("设置结束时间", DashboardTypeEndTime),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("设置间隔", DashboardTypeInterval),
			tgbotapi.NewInlineKeyboardButtonData("设置活动文案", DashboardTypeActivityText),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("设置发码数量", DashboardTypeCodeNum),
			tgbotapi.NewInlineKeyboardButtonData("开始任务", DashboardTypeStartTask),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("取消", cst.KeyboardTypeCancel),
		),
	)
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
