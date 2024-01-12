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
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/gin-gonic/gin"
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

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("COMMUNITY_")
	viper.AddConfigPath(*configFile)

	err = viper.ReadInConfig()
	utils.MustNotError(err)

	readableName := viper.GetString("readableName")
	err = initall.InitCommunity("COMMUNITY_")
	utils.MustNotError(err)

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
