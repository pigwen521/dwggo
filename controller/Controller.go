package controller

import (
	"net/http"

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
