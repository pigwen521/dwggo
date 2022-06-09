package consts

import "dsjk.com/dwggo/system/core"

//system用到的变量
const (
	REPORT_WARNING_URL  = "" //告警 系统错误日志会给该地址发post请求 msg:错误内容，form:WARNING_NAME；设置空字符串不告警
	REPORT_WARNING_NAME = "XXX平台"
)

func Init() {
	core.REPORT_WARNING_URL = REPORT_WARNING_URL
	core.REPORT_WARNING_NAME = REPORT_WARNING_NAME
}
