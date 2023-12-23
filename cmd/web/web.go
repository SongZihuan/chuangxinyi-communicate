package main

import (
	"flag"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth/login"
	initall "gitee.com/wuntsong/chuangxinyi-communicate/init"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/router"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var configFile = flag.String("f", "etc", "the config path")

func main() {
	flag.Parse()
	var err error

	//1.Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(0))

	//2.Set up configuration
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("COMMUNITY_")
	viper.AddConfigPath(*configFile)

	err = viper.ReadInConfig()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Parse conf file fail: %s", err.Error()))
		return
	}

	serviceName := viper.GetString("serviceName")
	err = initall.Init("COMMUNITY_", serviceName)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Fail to init: %s", err.Error()))
		return
	}

	_, err = login.ConnectWebSocket()
	if err != nil {
		logger.Logger.Error("websocket connect resp: %s", err)
	}

	engine := gin.Default()
	router.Setup(engine)
	_ = engine.Run(":" + viper.GetString("base.port"))
}
