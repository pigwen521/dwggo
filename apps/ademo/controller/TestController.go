package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**注释
 */
type TestController struct {
	ControllerBase
}

func (base *TestController) CallBefore(ctx *gin.Context, current_ctrl, current_action string) bool {
	return true
}
func (ctrl *TestController) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "test..."})
}
