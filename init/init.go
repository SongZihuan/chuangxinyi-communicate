package init

import (
	"context"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth"
	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/cron"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/global/peername"
	"gitee.com/wuntsong/chuangxinyi-communicate/ip"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/rand"
	"gitee.com/wuntsong/chuangxinyi-communicate/redis"
	"gitee.com/wuntsong/chuangxinyi-communicate/signalexit"
	"gitee.com/wuntsong/chuangxinyi-communicate/urls"
	"gitee.com/wuntsong/chuangxinyi-communicate/yundun"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

func InitCommunity(envPrefix string, serviceName string) errors.WTError {
	var err error
	err = signalexit.InitSignalExit(0)
	if err != nil {
		return errors.Errorf("signal resp: %s", errors.WarpQuick(err).Error())
	}

	signalexit.AddExitByFunc(func(ctx context.Context, _ os.Signal) context.Context {
		CloseAll()
		return context.WithValue(ctx, "InitClose", true)
	})

	mode := viper.GetString("mode")
	gin.SetMode(mode)

	err = logger.InitLogger(serviceName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = peername.InitPeerName(envPrefix)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = redis.InitRedis()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = dao.Setup()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = cache.Setup()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = cron.Setup()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = rand.InitRander()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = yundun.InitYunDun()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = auth.InitAuth()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = ip.InitYunIP()
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = urls.InitUrls()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
