package service

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"github.com/spf13/viper"
)

var AppinfoService = newAppinfoService()

func newAppinfoService() *appinfoService {
	return &appinfoService{}
}

type appinfoService struct {
}

func (s *appinfoService) GetAppinfo() *model.AppData {
	return &model.AppData{
		Name:           viper.GetString("readableName"),
		DomainID:       viper.GetString("auth.domainUID"),
		UserLevelAdmin: model.UserLevelAdmin,
	}
}
