package cache

import (
	"time"

	"github.com/goburrow/cache"

	"github.com/SongZihuan/chuangxinyi-communicate/dao"
)

type articleTagCache struct {
	cache cache.LoadingCache
}

var ArticleTagCache = newArticleTagCache()

func newArticleTagCache() *articleTagCache {
	return &articleTagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articleTags := dao.ArticleTagDao.FindByArticleId(key2Int64(key))
				if len(articleTags) > 0 {
					var tagIds []int64
					for _, articleTag := range articleTags {
						tagIds = append(tagIds, articleTag.TagId)
					}
					value = tagIds
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *articleTagCache) Get(articleId int64) []int64 {
	articleTags := dao.ArticleTagDao.FindByArticleId(articleId)
	if len(articleTags) > 0 {
		var tagIds []int64
		for _, articleTag := range articleTags {
			tagIds = append(tagIds, articleTag.TagId)
		}
		return tagIds
	}
	return nil
}

func (c *articleTagCache) Invalidate(articleId int64) {
	// 啥也不做
}
