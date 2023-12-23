package admin

import (
	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/controller"
	"gitee.com/wuntsong/chuangxinyi-communicate/util"
	// "gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

// SettingController setting controller
type SettingController struct {
	controller.BaseController
}

// List list settings
func (c *SettingController) List(ctx *gin.Context) {
	c.Success(ctx, service.SettingService.GetSetting())
}

// Store store settings
func (c *SettingController) Store(ctx *gin.Context) {
	config := ctx.Request.FormValue("config")
	if err := service.SettingService.SetAll(config); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}
