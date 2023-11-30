package utils

import (
	"strings"
)

// StringToStringSli :  字符串转切片
func StringToStringSli(str string) (sli []string) {
	if str == "" {
		return []string{}
	}
	delimiters := []string{",", ";", ":", "|", "-", "\n", "\t"}

	sign := ""
	for _, delimiter := range delimiters {
		parts := strings.Split(str, delimiter)
		if len(parts) > 1 {
			sign = delimiter
			break
		}
	}
	if sign == "" {
		sli = append(sli, str)
		return
	}
	strSli := strings.Split(str, sign)
	if len(strSli) > 0 {
		for _, v := range strSli {
			if v == "" {
				continue
			}
			sli = append(sli, v)
		}
	}
	return
}

// UniqueStrings 返回一个新的切片，其中包含原切片中的唯一字符串
func UniqueStrings(strings []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, str := range strings {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	return result
}
