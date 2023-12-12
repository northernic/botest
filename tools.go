package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-yaml/yaml"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// isValidInterval 可以是一个验证 cron 表达式有效性的函数
func isValidInterval(interval string) bool {
	pattern := `^(\d+h)?(\d+m)?(\d+s)?$`
	matched, _ := regexp.MatchString(pattern, interval)
	return matched

}

func checkAuth(groupID int64, moduleName string) bool {
	//先查找groupID权限
	authID := ""
	for _, v := range globalConf.GroupAuth {
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
	moduleName = strings.ToLower(moduleName)
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
	case "aquis":
		return Aquis
	case "voya":
		return Voya
	case "jinsha":
		return JinSha
	case "shangpujing":
		return ShangPuJing
	case "jason":
		return Jason
	default:
		return 0
	}

}

// mentionBot检查消息是否包含对机器人的提及
func mentionBot(message *tgbotapi.Message, botUsername string) bool {
	if message == nil {
		return false
	}
	if message.Entities == nil {
		return false
	}
	for _, entity := range *message.Entities {
		if entity.Type == "mention" {
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			if mention == "@"+botUsername {
				return true
			}
		}
	}
	return false
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

func formatCronEntries(entries []cron.Entry) string {
	var sb strings.Builder
	for _, entry := range entries {
		sb.WriteString(fmt.Sprintf("Entry ID: %d\n", entry.ID))
		sb.WriteString(fmt.Sprintf("Next Run: %s\n", entry.Next))
		sb.WriteString(fmt.Sprintf("Prev Run: %s\n", entry.Prev))
		sb.WriteString("----------\n")
	}
	return sb.String()
}
