package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	ControllerBase
}

func (ctrl *HomeController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home/index.html", gin.H{"title": "welcome!"})
}
