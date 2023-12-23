package main

import (
	"flag"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth/login"
	initall "gitee.com/wuntsong/chuangxinyi-communicate/init"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/router"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/log"
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
		log.Fatal(fmt.Sprintf("Parse conf file fail: %s", err.Error()))
	}

	serviceName := viper.GetString("serviceName")
	err = initall.Init("COMMUNITY_", serviceName)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to init: %s", err.Error()))
	}

	_, err = login.ConnectWebSocket()
	if err != nil {
		logger.Logger.Error("websocket connect resp: %s", err)
	}

	engine := gin.Default()
	router.Setup(engine)
	_ = engine.Run(":" + viper.GetString("base.port"))
}
