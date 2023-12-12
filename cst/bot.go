package cst

import "strings"

const (
	ChatTypeGroup string = "group"
)

type CmdType int

const (
	CmdTypeNone   CmdType = iota + 1 // 非命令
	CmdTypeSingle                    // 单命令
	CmdTypeMul                       //多重命令
)

func GetCmdType(cmd string) CmdType {
	if cmd == "" {
		return CmdTypeNone
	}
	cmds := strings.Split(cmd, "/")
	if len(cmds) == 2 && cmds[0] == "" && cmds[1] != "" {
		return CmdTypeSingle
	}

	if len(cmds) >= 3 && cmds[1] != "" && cmds[2] != "" {
		return CmdTypeMul
	}
	return CmdTypeNone
}

// 键盘类型
const (
	KeyboardTypeMain        string = "主菜单"      //主菜单
	KeyboardTypeSettingTask string = "设定定时任务菜单" //设定定时任务菜单
	KeyboardTypeTaskList    string = "定时任务菜单"   //定时任务菜单
	KeyboardTypeCancel      string = "取消"       //取消
)
