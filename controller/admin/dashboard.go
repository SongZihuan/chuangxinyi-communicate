package admin

import (
	"github.com/gin-gonic/gin"
	"runtime"
	"time"

	"gitee.com/wuntsong/chuangxinyi-communicate/config"
	"gitee.com/wuntsong/chuangxinyi-communicate/controller"
	"gitee.com/wuntsong/chuangxinyi-communicate/util"
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
		"appName":   config.AppName,
		"upTime":    util.TimeSincePro(initTime),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"numCpu":    runtime.NumCPU(),
		"goversion": runtime.Version(),
	})
}
