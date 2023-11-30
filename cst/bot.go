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
