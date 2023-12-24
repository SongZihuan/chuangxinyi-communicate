package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"runtime"
	"time"

	"gitee.com/wuntsong/chuangxinyi-communicate/controller"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

// initTime is the time when the application was initialized.
var initTime = time.Now()

// DashboardController dashboard controller
type DashboardController struct {
	controller.BaseController
}

// GetSysteminfo get system info
func (c *DashboardController) Systeminfo(ctx *gin.Context) {
	c.Success(ctx, gin.H{
		"appName":   viper.GetString("readableName"),
		"upTime":    utils.TimeSincePro(initTime),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"numCpu":    runtime.NumCPU(),
		"goversion": runtime.Version(),
	})
}
