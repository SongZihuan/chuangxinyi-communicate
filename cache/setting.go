package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/redis"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

type settingCache struct {
}

var SettingCache = newSettingCache()

func newSettingCache() *settingCache {
	return &settingCache{}
}

func (c *settingCache) Get(key string) *model.Setting {
	redisKey := fmt.Sprintf("settings:%s", key)
	dataString, ok := redis.GetCache(context.Background(), redisKey)
	if ok {
		var data model.Setting
		err := utils.JsonUnmarshal([]byte(dataString), &data)
		if err == nil {
			return &data
		}
	}

	val := dao.SettingDao.GetByKey(key)
	valString, err := json.Marshal(val)
	if err == nil {
		redis.SetCache(context.Background(), redisKey, string(valString), 0)
	}

	return val
}

func (c *settingCache) GetValue(key string) string {
	data, ok := redis.GetCache(context.Background(), fmt.Sprintf("settings:%s", key))
	if ok {
		return data
	}
	return ""
}

func (c *settingCache) Invalidate(key string) {
	redis.DelCache(context.Background(), fmt.Sprintf("settings:%s", key))
}
