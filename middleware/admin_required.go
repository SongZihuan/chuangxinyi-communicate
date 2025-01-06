package middleware

import (
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminRequired admin required
func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.UserService.GetCurrent(ctx)
		if user == nil {
			err := utils.ErrorNotLogin
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}
		if user.Level != model.UserLevelAdmin {
			err := utils.ErrorPermissionDenied
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}
		ctx.Next()
	}
}
