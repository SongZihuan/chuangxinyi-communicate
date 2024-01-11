package main

import (
	"context"
	"flag"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth/login"
	initall "gitee.com/wuntsong/chuangxinyi-communicate/init"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/router"
	"gitee.com/wuntsong/chuangxinyi-communicate/signalexit"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var configFile = flag.String("f", "etc", "the config path")

func main() {
	fmt.Println("Start community backend...")
	CmdMain()
	fmt.Println("Bye~")
}

func CmdMain() {
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
	readableName := viper.GetString("readableName")
	err = initall.InitCommunity("COMMUNITY_", serviceName)
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

	addr := ":" + viper.GetString("base.port")
	server := http.Server{
		Addr:    addr,
		Handler: NewServer(engine),
	}

	signalexit.AddExitByFunc(func(ctx context.Context, signal os.Signal) context.Context {
		_ = server.Shutdown(ctx)
		return context.WithValue(ctx, "Server-Shutdown", true)
	})

	logger.Logger.WXInfo("启动服务 %s 在端口 %s...", readableName, addr)
	go func() {
		_ = server.ListenAndServe()
	}()

	select {} // 阻塞
}
