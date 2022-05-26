package core

/*
lishaowen
409250643@qq.com
*/
import (
	"net/http"
	"path/filepath"

	"dsjk.com/openplatform/consts"
	"github.com/gin-gonic/gin"
)

func init() {
	initSysVars()
}
func initSysVars() {
	var err error
	consts.APP_ROOT, err = filepath.Abs("")
	if err != nil {
		panic("初始化失败：" + err.Error())
	}
	consts.APP_ROOT += "/"
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
