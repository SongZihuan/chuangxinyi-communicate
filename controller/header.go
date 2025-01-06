package controller

import (
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/gin-gonic/gin"
)

type HeaderController struct {
	BaseController
}

func (c *HeaderController) GetHeader(ctx *gin.Context) {
	uid := ctx.Request.FormValue("uid")
	url := service.UserService.GetHeader(uid)
	c.Redirect(ctx, url)
}
