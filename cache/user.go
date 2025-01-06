package cache

import (
	"time"

	"github.com/goburrow/cache"

	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

type userCache struct {
	cache      cache.LoadingCache
	scoreCache cache.LoadingCache
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, err error) {
				value = dao.UserDao.Get(key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		scoreCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, err error) {
				userScore := dao.UserScoreDao.FindOne(sqlcnd.NewSqlCnd().Eq("user_id", key2Int64(key)))
				if userScore == nil {
					value = 0
				} else {
					value = userScore.Score
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	return dao.UserDao.Get(userId)
}

func (c *userCache) Invalidate(userId int64) {
	// 啥也不做
}

func (c *userCache) GetScore(userId int64) int {
	userScore := dao.UserScoreDao.FindOne(sqlcnd.NewSqlCnd().Eq("user_id", userId))
	if userScore == nil {
		return 0
	} else {
		return userScore.Score
	}
}

func (c *userCache) InvalidateScore(userId int64) {
	// 啥也不做
}
