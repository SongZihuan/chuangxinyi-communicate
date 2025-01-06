package init

import (
	"context"
	"github.com/SongZihuan/chuangxinyi-communicate/auth"
	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/cron"
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/global/peername"
	"github.com/SongZihuan/chuangxinyi-communicate/ip"
	"github.com/SongZihuan/chuangxinyi-communicate/logger"
	"github.com/SongZihuan/chuangxinyi-communicate/rand"
	"github.com/SongZihuan/chuangxinyi-communicate/redis"
	"github.com/SongZihuan/chuangxinyi-communicate/signalexit"
	"github.com/SongZihuan/chuangxinyi-communicate/urls"
	"github.com/SongZihuan/chuangxinyi-communicate/yundun"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"os"
)

func InitCommunity(envPrefix string) errors.WTError {
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

	err = logger.InitLogger(viper.GetString("serviceName"))
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

func InitSqlClear(envPrefix string) errors.WTError {
	var err error
	err = signalexit.InitSignalExit(0)
	if err != nil {
		return errors.Errorf("signal resp: %s", errors.WarpQuick(err).Error())
	}

	signalexit.AddExitByFunc(func(ctx context.Context, _ os.Signal) context.Context {
		CloseAll()
		return context.WithValue(ctx, "InitClose", true)
	})

	err = logger.InitLogger(viper.GetString("serviceName"))
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

	err = rand.InitRander()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
