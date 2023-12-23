package main

import (
	"flag"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/cron"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/middleware"
	"gitee.com/wuntsong/chuangxinyi-communicate/router"
	"gitee.com/wuntsong/chuangxinyi-communicate/util/log"
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
	viper.AddConfigPath(*configFile)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse conf file fail: %s", err.Error()))
	}

	//3.Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)

	//4.Set up database connection
	dao.Setup()

	//5.Set up cache
	cache.Setup()

	////5.Set up cron
	cron.Setup()

	//6.Initialize language
	middleware.InitLang()

	engine := gin.Default()
	router.Setup(engine)
	_ = engine.Run(":" + viper.GetString("base.port"))
}
