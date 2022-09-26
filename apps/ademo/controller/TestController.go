package controller

import (
	"net/http"
	"time"

	"dsjk.com/dwggo/system/core"
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
func (ctrl *TestController) LogError(ctx *gin.Context) {
	core.LogError("test" + time.Now().String())
	ctx.JSON(http.StatusOK, gin.H{"msg": "logerror..."})
}
