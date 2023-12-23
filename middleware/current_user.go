package middleware

import (
	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

func CurrentUser(ctx *gin.Context) {
	ctx.Set("CurrentUser", service.UserService.GetCurrent(ctx))
	ctx.Next()
}
