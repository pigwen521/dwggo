package controller

import (
	"net/http"

	"dsjk.com/dwggo/apps/ademo/model"
	"dsjk.com/dwggo/apps/ademo/service"
	"dsjk.com/dwggo/system/core"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

/**
XXX控制器
*/
type AdemoController struct {
	ControllerBase
}

/**
钩子方法，在action运行之前执行
重写该方法，return false可阻断action的执行
*/
func (base *AdemoController) CallBefore(ctx *gin.Context, current_ctrl, current_action string) bool {
	return true
}
func (ctrl *AdemoController) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello index")
}
func (ctrl *AdemoController) Query(ctx *gin.Context) {
	var arg_in model.ArgAdemoQueryInModel
	var ademoService service.AdemoService

	ctx.ShouldBindWith(&arg_in, binding.Query)

	arg_out, err := ademoService.Query(arg_in)
	if err != nil {
		core.ResultFail(ctx, err.Error())
	} else {
		core.ResultSucc(ctx, arg_out)
	}
}
