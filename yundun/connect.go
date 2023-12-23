package yundun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/viper"
)

var YunDunClient *green20220302.Client

func InitYunDun() error {
	var err error

	accessKeyId := viper.GetString("aliyun.accessKeyId")
	accessKeySecret := viper.GetString("aliyun.accessKeySecret")

	YunDunClient, err = green20220302.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("green-cip.cn-shanghai.aliyuncs.com"),
	})
	if err != nil {
		return err
	}

	return nil
}
