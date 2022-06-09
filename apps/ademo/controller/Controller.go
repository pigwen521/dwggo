package controller

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"dsjk.com/dwggo/system/core"
	"github.com/gin-gonic/gin"
)

type Controller interface {
}
type ControllerBase struct {
}

func (base *ControllerBase) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello index")
}

/**
钩子方法，在action运行之前执行
重写该方法，return false可阻断action的执行
*/
func (base *ControllerBase) CallBefore(ctx *gin.Context, current_action string) bool {
	return true
}

func InitCtrlByName(str string) (reflect.Value, error) {
	var v reflect.Value
	switch strings.ToLower(str) {
	case "ademo":
		v = reflect.ValueOf(new(AdemoController))
	case "home":
		v = reflect.ValueOf(new(HomeController))
	case "test":
		v = reflect.ValueOf(new(TestController))
	//新增相应控制器代码
	default:
		if core.IsEnvDev() {
			return reflect.Value{}, errors.New("控制器未配置？请在controller/Controller.go InitCtrlByName方法中完善case:" + str)
		} else {
			return reflect.Value{}, errors.New("this controller is wrong：" + str)
		}
	}
	return v, nil
}
