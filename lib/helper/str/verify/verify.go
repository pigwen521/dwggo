package verify

import (
	"regexp"
)

//格式验证

//日期 YYYY-DD-MM
func IsDate(str *string) bool {
	m, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, *str)
	return m
}

//时间 13:30
func IsTime(str *string) bool {
	m, _ := regexp.MatchString(`^\d{2}:\d{2}$`, *str)
	return m
}

//纯数字
func IsNumber(str *string) bool {
	m, _ := regexp.MatchString(`^\d+$`, *str)
	return m
}
