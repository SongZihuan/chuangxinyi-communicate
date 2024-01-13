package cache

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
)

var (
	topicRecommendCacheKey = "recommend_topics_cache"
)

var TopicCache = newTopicCache()

type topicCache struct {
}

func newTopicCache() *topicCache {
	return &topicCache{}
}

func (c *topicCache) GetRecommendTopics() []model.Topic {
	return dao.TopicDao.Find(sqlcnd.NewSqlCnd().Eq("recommend", true).Eq("status", model.StatusOk).Limit(20).Desc("last_comment_time"))
}

func (c *topicCache) InvalidateRecommend() {
	// 啥也不做
}
