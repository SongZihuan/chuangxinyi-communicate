package middleware

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SigninRequired signin required
func SigninRequired(ctx *gin.Context) {
	user := service.UserService.GetCurrent(ctx)
	if user == nil {
		err := utils.ErrorNotLogin
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    err.Code,
			"message": err.Message,
		})

	}
	ctx.Next()
}
