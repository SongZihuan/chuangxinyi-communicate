package controller

import (
	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

type ConfigController struct {
	BaseController
}

func (c *ConfigController) List(ctx *gin.Context) {

	data := map[string]interface{}{}
	data["setting"] = service.SettingService.GetSetting()
	data["appinfo"] = service.AppinfoService.GetAppinfo()

	c.Success(ctx, data)
}
