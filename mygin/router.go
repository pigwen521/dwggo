package mygin

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"dsjk.com/dwggo/controller"
	"dsjk.com/dwggo/core"
	"dsjk.com/dwggo/lib/helper/str"

	"github.com/gin-gonic/gin"
)

func InitCtrlByName(str string) (reflect.Value, error) {
	var v reflect.Value
	switch strings.ToLower(str) {
	case "ademo":
		v = reflect.ValueOf(new(controller.AdemoController))
	//待增加
	default:
		if core.IsEnvDev() {
			return reflect.Value{}, errors.New("控制器未配置？请在mygin/router.go InitCtrlByName方法中完善case:" + str)
		} else {
			return reflect.Value{}, errors.New("this controller is wrong：" + str)
		}
	}
	return v, nil
}

//通过URL中的action，找出对应的方法
func callMethodByAction(v reflect.Value, act string, ctx *gin.Context) ([]reflect.Value, error) {
	k := v.Type()
	for i := 0; i < v.NumMethod(); i++ {
		key := k.Method(i)
		//val := v.Method(i)
		//fmt.Println("getMethodByAction", i, key.Name, val.Type(), val.Interface())
		if strings.ToLower(key.Name) == strings.ToLower(act) { //URL中的大小写和action不一致
			params := make([]reflect.Value, 2)
			params[0] = v
			params[1] = reflect.ValueOf(ctx)
			return key.Func.Call(params), nil //v.Method(i).Call()
		}
	}
	return nil, errors.New("have not find the action to call:" + act)

	/*
		方式2，执行反射的方法MethodByName
		v := InitCtrlByName(ctrl)

		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(ctx)
		v.MethodByName(action).Call(params) */
}
func callMiddleWare(v reflect.Value, action string, ctx *gin.Context, cur_action string) []reflect.Value {
	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(ctx)
	params[1] = reflect.ValueOf(cur_action)
	return v.MethodByName(action).Call(params)
}
func PageNotFound(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error/404.html", gin.H{"title": "404"})
}

func InitRoute(r *gin.Engine) {
	root_path := core.GetConfigString("app.path")
	if root_path != "/" { //站点目录非根目录
		r.GET("/", func(ctx *gin.Context) {
			routerZeroLevel(ctx)
		}).POST("/", func(ctx *gin.Context) {
			routerZeroLevel(ctx)
		})
	}

	r.GET(root_path, func(ctx *gin.Context) {
		routerZeroLevel(ctx)
	}).POST(root_path, func(ctx *gin.Context) {
		routerZeroLevel(ctx)
	})

	r.GET(root_path+":ctrl", func(ctx *gin.Context) {
		routerOneLevel(ctx)
	}).POST(root_path+":ctrl", func(ctx *gin.Context) {
		routerOneLevel(ctx)
	})

	r.GET(root_path+":ctrl/:action", func(ctx *gin.Context) {
		routerTwoLevel(ctx)
	}).POST(root_path+":ctrl/:action", func(ctx *gin.Context) {
		routerTwoLevel(ctx)
	})
}

//:ctrl/:action 两级路由
func routerTwoLevel(ctx *gin.Context) {
	ctrl := ctx.Param("ctrl")
	action := str.FirstToUpper(ctx.Param("action"))
	executCtrlAction(ctx, ctrl, action)
}

//:ctrl 没action 一级路由
func routerOneLevel(ctx *gin.Context) {
	ctrl := ctx.Param("ctrl")
	action := core.GetConfigString("router.default_action")
	executCtrlAction(ctx, ctrl, action)
}

//根目录首页 0级路由
func routerZeroLevel(ctx *gin.Context) {
	ctrl := core.GetConfigString("router.default_controller")
	action := core.GetConfigString("router.default_action")
	executCtrlAction(ctx, ctrl, action)
}

/**
执行控制器的方法
*/
func executCtrlAction(ctx *gin.Context, ctrl string, action string) {
	defer core.ForPanicLog()
	v, err := InitCtrlByName(ctrl)
	if err != nil {
		core.LogWarning(err.Error())
		PageNotFound(ctx)
		return
	}

	call_ret := callMiddleWare(v, "CallBefore", ctx, action)
	if !call_ret[0].Interface().(bool) {
		return
	}

	_, err = callMethodByAction(v, action, ctx)
	if err != nil {
		core.LogWarning(err.Error())
		PageNotFound(ctx)
		return
	}
}
