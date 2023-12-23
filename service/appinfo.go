package service

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/config"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
)

var AppinfoService = newAppinfoService()

func newAppinfoService() *appinfoService {
	return &appinfoService{}
}

type appinfoService struct {
}

func (s *appinfoService) GetAppinfo() *model.AppData {
	return &model.AppData{
		Name:           config.AppName,
		UserLevelAdmin: model.UserLevelAdmin,
	}
}
