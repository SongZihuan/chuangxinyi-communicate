package main

import (
	"flag"
	"fmt"
	initall "gitee.com/wuntsong/chuangxinyi-communicate/init"
	"gitee.com/wuntsong/chuangxinyi-communicate/sqlclear"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/spf13/viper"
)

var configFile = flag.String("f", "etc", "the config path")

func main() {
	fmt.Println("Start sql clear...")
	CmdMain()
	fmt.Println("Bye~")
}

func CmdMain() {
	var err error
	flag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("COMMUNITY_")
	viper.AddConfigPath(*configFile)

	err = viper.ReadInConfig()
	utils.MustNotError(err)

	err = initall.InitSqlClear("COMMUNITY_")
	utils.MustNotError(err)

	err = sqlclear.StartClear()
	utils.MustNotError(err)
}
