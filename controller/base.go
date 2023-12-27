package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

// BaseController controller
type BaseController struct {
}

// BindAndValidate bind and validate
func (c *BaseController) BindAndValidate(ctx *gin.Context, obj interface{}) bool {
	if err := form.Bind(ctx, obj); err != nil {
		c.Fail(ctx, utils.NewError(-1, err.Error()))
		return false
	}
	return true
}

// GetCurrentUser get current user from contexg
func (c *BaseController) GetCurrentUser(ctx *gin.Context) *model.User {
	if currentUser, ok := ctx.Get("CurrentUser"); ok {
		return currentUser.(*model.User)
	}
	return nil
}

// Success output json data
func (c *BaseController) Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
		"success": true,
	})
}

// Redirect redirect to url
func (c *BaseController) Redirect(ctx *gin.Context, url string) {
	ctx.Redirect(http.StatusFound, url)
}

// Fail output error
func (c *BaseController) Fail(ctx *gin.Context, error *utils.CodeError) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    error.CodeInt(),
		"message": error.Msg(),
		"success": false,
	})
	return
}
