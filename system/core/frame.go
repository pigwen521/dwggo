package core

/*
lishaowen
409250643@qq.com
*/
import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//系统初始化的一些全局变量
var APP_ROOT string            //项目根目录，main.go的目录
var APILOG_FILENAME string     //API请求日志文件名
var REPORT_WARNING_URL string  //告警URL
var REPORT_WARNING_NAME string //告警来源名称

func init() {
	initSysVars()
}
func initSysVars() {
	var err error
	APP_ROOT, err = filepath.Abs("")
	if err != nil {
		panic("初始化失败：" + err.Error())
	}
	APP_ROOT += "/"
	APILOG_FILENAME = GetConfigString("logger.apilog_filename")
}

const (
	CODE_SUCC = 1
	CODE_FAIL = -1
)

func ResultFail(ctx *gin.Context, msg string) {
	ReslutJson(ctx, CODE_FAIL, nil, msg)
}
func ResultSucc(ctx *gin.Context, data interface{}) {
	ReslutJson(ctx, CODE_SUCC, data, "")
}
func ReslutJson(ctx *gin.Context, code int, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, ReturnJson(ctx, code, data, msg))
}

func ReturnFail(ctx *gin.Context, msg string) interface{} {
	return ReturnJson(ctx, CODE_FAIL, nil, msg)
}
func ReturnSucc(ctx *gin.Context, data interface{}) interface{} {
	return ReturnJson(ctx, CODE_SUCC, data, "")
}
func ReturnJson(ctx *gin.Context, code int, data interface{}, msg string) interface{} {
	return gin.H{"code": code, "data": data, "msg": msg}
}

func SiteUrl(uri string) string {
	return GetConfigString("app.site_root") + uri
}

/**
返回当前请求的actions
./user/act1/act2
return [act1,act2]
*/
func GetActions(ctx *gin.Context) []string {
	actions := strings.Trim(ctx.Param("actions"), "/") // / /act1	/act1/act2	/act1/act2/act3
	if actions == "" {
		actions = strings.ToLower(GetConfigString("router.default_action"))
	}
	return strings.Split(actions, "/")
}
