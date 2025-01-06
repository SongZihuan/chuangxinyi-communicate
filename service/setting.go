package service

import (
	errors "github.com/wuntsong-org/wterrors"
	"strconv"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/logger"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

var SettingService = newSettingService()

func newSettingService() *settingService {
	return &settingService{}
}

type settingService struct {
}

func (s *settingService) Get(id int64) *model.Setting {
	return dao.SettingDao.Get(id)
}

func (s *settingService) Take(where ...interface{}) *model.Setting {
	return dao.SettingDao.Take(where...)
}

func (s *settingService) Find(cnd *sqlcnd.SqlCnd) []model.Setting {
	return dao.SettingDao.Find(cnd)
}

func (s *settingService) FindOne(cnd *sqlcnd.SqlCnd) *model.Setting {
	return dao.SettingDao.FindOne(cnd)
}

func (s *settingService) List(cnd *sqlcnd.SqlCnd) (list []model.Setting, paging *sqlcnd.Paging) {
	return dao.SettingDao.List(cnd)
}

func (s *settingService) GetAll() []model.Setting {
	return dao.SettingDao.Find(sqlcnd.NewSqlCnd().Asc("id"))
}

func (s *settingService) SetAll(configStr string) errors.WTError {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return dao.Tx(dao.DB(), func(tx *gorm.DB) errors.WTError {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return errors.WarpQuick(err)
			}
		}
		return nil
	})
}

// 设置配置，如果配置不存在，那么创建
func (s *settingService) Set(key, value, name, description string) errors.WTError {
	return dao.Tx(dao.DB(), func(tx *gorm.DB) errors.WTError {
		if err := s.setSingle(tx, key, value, name, description); err != nil {
			return errors.WarpQuick(err)
		}
		return nil
	})
}

func (s *settingService) setSingle(db *gorm.DB, key, value, name, description string) errors.WTError {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := dao.SettingDao.GetByKey(key)
	if sysConfig == nil {
		sysConfig = &model.Setting{
			CreateTime: utils.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = utils.NowTimestamp()

	if len(name) > 0 {
		sysConfig.Name = name
	}
	if len(description) > 0 {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.ID > 0 {
		err = dao.SettingDao.Update(sysConfig)
	} else {
		err = dao.SettingDao.Create(sysConfig)
	}
	if err != nil {
		return errors.WarpQuick(err)
	} else {
		cache.SettingCache.Invalidate(key)
		return nil
	}
}

func (s *settingService) GetSetting() *model.ConfigData {
	var (
		siteTitle        = cache.SettingCache.GetValue(model.SettingSiteTitle)
		siteDescription  = cache.SettingCache.GetValue(model.SettingSiteDescription)
		siteKeywords     = cache.SettingCache.GetValue(model.SettingSiteKeywords)
		siteNavs         = cache.SettingCache.GetValue(model.SettingSiteNavs)
		siteTips         = cache.SettingCache.GetValue(model.SettingSiteTips)
		siteNotification = cache.SettingCache.GetValue(model.SettingSiteNotification)
		recommendTags    = cache.SettingCache.GetValue(model.SettingRecommendTags)
		scoreConfigStr   = cache.SettingCache.GetValue(model.SettingScoreConfig)
		footerConfigStr  = cache.SettingCache.GetValue(model.SettingFooterConfig)
		defaultNodeIdStr = cache.SettingCache.GetValue(model.SettingDefaultNodeId)
	)

	var siteKeywordsArr []string
	if err := utils.ParseJson(siteKeywords, &siteKeywordsArr); err != nil {
		logger.Logger.Info("站点关键词数据错误")
	}

	var siteNavsArr []model.SiteNav
	if err := utils.ParseJson(siteNavs, &siteNavsArr); err != nil {
		logger.Logger.Info("站点导航数据错误")
	}

	var siteTipsArr []model.SiteTip
	if err := utils.ParseJson(siteTips, &siteTipsArr); err != nil {
		logger.Logger.Info("小贴士数据错误")
	}

	var recommendTagsArr []string
	if err := utils.ParseJson(recommendTags, &recommendTagsArr); err != nil {
		logger.Logger.Info("推荐标签数据错误")
	}

	var scoreConfig model.ScoreConfig
	if err := utils.ParseJson(scoreConfigStr, &scoreConfig); err != nil {
		logger.Logger.Info("积分配置错误")
	}

	var footerConfig model.FooterConfig
	if err := utils.ParseJson(footerConfigStr, &footerConfig); err != nil {
		logger.Logger.Info("底部信息配置错误")
	}

	var defaultNodeId, _ = strconv.ParseInt(defaultNodeIdStr, 10, 64)

	return &model.ConfigData{
		SiteTitle:        siteTitle,
		SiteDescription:  siteDescription,
		SiteKeywords:     siteKeywordsArr,
		SiteNavs:         siteNavsArr,
		SiteTips:         siteTipsArr,
		SiteNotification: siteNotification,
		RecommendTags:    recommendTagsArr,
		ScoreConfig:      scoreConfig,
		FooterConfig:     footerConfig,
		DefaultNodeId:    defaultNodeId,
	}
}
