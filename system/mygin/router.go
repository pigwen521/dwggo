package mygin

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper/str"

	"github.com/gin-gonic/gin"
)

type InitCtrlByNameCB func(str string) (reflect.Value, error)

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
func callMiddleWare(v reflect.Value, action string, ctx *gin.Context, cur_ctrl, cur_action string) []reflect.Value {
	params := make([]reflect.Value, 3)
	params[0] = reflect.ValueOf(ctx)
	params[1] = reflect.ValueOf(strings.ToLower(cur_ctrl))
	params[2] = reflect.ValueOf(strings.ToLower(cur_action))
	return v.MethodByName(action).Call(params)
}
func PageNotFound(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error/404.html", gin.H{"title": "404"})
}

func InitRoute(r *gin.Engine, initCtrlByNameCB InitCtrlByNameCB) {
	root_path := core.GetConfigString("app.path")

	if root_path != "/" { //站点目录非根目录
		r.GET("/", func(ctx *gin.Context) {
			routerZeroLevel(ctx, initCtrlByNameCB)
		}).POST("/", func(ctx *gin.Context) {
			routerZeroLevel(ctx, initCtrlByNameCB)
		})
	}

	r.GET(root_path, func(ctx *gin.Context) {
		routerZeroLevel(ctx, initCtrlByNameCB)
	}).POST(root_path, func(ctx *gin.Context) {
		routerZeroLevel(ctx, initCtrlByNameCB)
	})

	r.GET(root_path+":ctrl/*actions", func(ctx *gin.Context) {
		routerAllLevel(ctx, initCtrlByNameCB)
	}).POST(root_path+":ctrl/*actions", func(ctx *gin.Context) {
		routerAllLevel(ctx, initCtrlByNameCB)
	})

}

//:ctrl/*actions 全部。多级
func routerAllLevel(ctx *gin.Context, initCtrlByNameCB InitCtrlByNameCB) {
	ctrl := ctx.Param("ctrl")
	action := ""
	actions := ctx.Param("actions") // / /act1	/act1/act2	/act1/act2/act3
	act_arr := strings.Split(actions, "/")
	if actions == "/" {
		action = str.FirstToUpper(core.GetConfigString("router.default_action"))
	} else {
		action = str.FirstToUpper(act_arr[1])
	}
	executCtrlAction(ctx, ctrl, action, initCtrlByNameCB)
}

//根目录首页 0级路由
func routerZeroLevel(ctx *gin.Context, initCtrlByNameCB InitCtrlByNameCB) {
	ctrl := core.GetConfigString("router.default_controller")
	action := core.GetConfigString("router.default_action")
	executCtrlAction(ctx, ctrl, action, initCtrlByNameCB)
}

/**
执行控制器的方法
*/
func executCtrlAction(ctx *gin.Context, ctrl string, action string, initCtrlByNameCB InitCtrlByNameCB) {
	defer core.ForPanicLog()
	v, err := initCtrlByNameCB(ctrl)
	if err != nil {
		core.LogWarning(err.Error())
		PageNotFound(ctx)
		return
	}
	call_ret := callMiddleWare(v, "CallBefore", ctx, ctrl, action)
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
