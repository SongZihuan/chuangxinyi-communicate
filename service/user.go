package service

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/auth/header"
	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/log"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ScanUserCallback func(users []model.User)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return dao.UserDao.Get(id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return dao.UserDao.Take(where...)
}

func (s *userService) Find(cnd *sqlcnd.SqlCnd) []model.User {
	return dao.UserDao.Find(cnd)
}

func (s *userService) FindOne(cnd *sqlcnd.SqlCnd) *model.User {
	return dao.UserDao.FindOne(cnd)
}

func (s *userService) List(cnd *sqlcnd.SqlCnd) (list []model.User, paging *sqlcnd.Paging) {
	return dao.UserDao.List(cnd)
}

func (s *userService) Count(cnd *sqlcnd.SqlCnd) int {
	return dao.UserDao.Count(cnd)
}

func (s *userService) Update(dto form.UserUpdateForm) error {
	err := dao.UserDao.Updates(dto.ID, map[string]interface{}{
		"level":       dto.Level,
		"update_time": utils.NowTimestamp(),
	})
	cache.UserCache.Invalidate(dto.ID)

	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := dao.UserDao.Updates(id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := dao.UserDao.UpdateColumn(id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	dao.UserDao.Delete(id)
	cache.UserCache.Invalidate(id)
}

// 获取当前登录用户
func (s *userService) GetCurrent(ctx *gin.Context) *model.User {
	userTmp, _ := ctx.Get(viper.GetString("jwt.identity_key"))
	if userTmp == nil {
		return nil
	}

	user := userTmp.(model.UserClaims).User
	if user == nil || user.Status != model.StatusOk {
		return nil
	}
	return user
}

// Scan 扫描
func (s *userService) Scan(cb ScanUserCallback) {
	var cursor int64
	for {
		list := dao.UserDao.Find(sqlcnd.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

// IncrTopicCount topic_count + 1
func (s *userService) IncrTopicCount(userId int64) int {
	t := dao.UserDao.Get(userId)
	if t == nil {
		return 0
	}
	topicCount := t.TopicCount + 1
	if err := dao.UserDao.UpdateColumn(userId, "topic_count", topicCount); err != nil {
		log.Error(err.Error())
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return topicCount
}

// IncrCommentCount comment_count + 1
func (s *userService) IncrCommentCount(userId int64) int {
	t := dao.UserDao.Get(userId)
	if t == nil {
		return 0
	}
	commentCount := t.CommentCount + 1
	if err := dao.UserDao.UpdateColumn(userId, "comment_count", commentCount); err != nil {
		log.Error(err.Error())
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return commentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := dao.TopicDao.Count(sqlcnd.NewSqlCnd().Eq("user_id", user.ID).Eq("status", model.StatusOk))
			commentCount := dao.CommentDao.Count(sqlcnd.NewSqlCnd().Eq("user_id", user.ID).Eq("status", model.StatusOk))
			_ = dao.UserDao.UpdateColumn(user.ID, "topic_count", topicCount)
			_ = dao.UserDao.UpdateColumn(user.ID, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.ID)
		}
	})
}

// GetHeader 获取头像url
func (s *userService) GetHeader(uid string) string {
	return header.GetHeaderUrl(uid)
}
