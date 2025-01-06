package service

import (
	errors "github.com/wuntsong-org/wterrors"

	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

var FavoriteService = newFavoriteService()

func newFavoriteService() *favoriteService {
	return &favoriteService{}
}

type favoriteService struct {
}

func (s *favoriteService) Get(id int64) *model.Favorite {
	return dao.FavoriteDao.Get(id)
}

func (s *favoriteService) Take(where ...interface{}) *model.Favorite {
	return dao.FavoriteDao.Take(where...)
}

func (s *favoriteService) Find(cnd *sqlcnd.SqlCnd) []model.Favorite {
	return dao.FavoriteDao.Find(cnd)
}

func (s *favoriteService) FindOne(cnd *sqlcnd.SqlCnd) *model.Favorite {
	return dao.FavoriteDao.FindOne(cnd)
}

func (s *favoriteService) List(cnd *sqlcnd.SqlCnd) (list []model.Favorite, paging *sqlcnd.Paging) {
	return dao.FavoriteDao.List(cnd)
}

func (s *favoriteService) Create(t *model.Favorite) errors.WTError {
	return dao.FavoriteDao.Create(t)
}

func (s *favoriteService) Update(t *model.Favorite) errors.WTError {
	return dao.FavoriteDao.Update(t)
}

func (s *favoriteService) Updates(id int64, columns map[string]interface{}) errors.WTError {
	return dao.FavoriteDao.Updates(id, columns)
}

func (s *favoriteService) UpdateColumn(id int64, name string, value interface{}) errors.WTError {
	return dao.FavoriteDao.UpdateColumn(id, name, value)
}

func (s *favoriteService) Delete(id int64) {
	dao.FavoriteDao.Delete(id)
}

func (s *favoriteService) GetBy(userId int64, entityType string, entityId int64) *model.Favorite {
	return dao.FavoriteDao.Take("user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId)
}

// 收藏文章
func (s *favoriteService) AddArticleFavorite(userId, articleId int64) errors.WTError {
	article := dao.ArticleDao.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return errors.New("收藏的文章不存在")
	}
	return s.addFavorite(userId, model.EntityTypeArticle, articleId)
}

// 收藏主题
func (s *favoriteService) AddTopicFavorite(userId, topicId int64) errors.WTError {
	topic := dao.TopicDao.Get(topicId)
	if topic == nil || topic.Status != model.StatusOk {
		return errors.New("收藏的话题不存在")
	}
	return s.addFavorite(userId, model.EntityTypeTopic, topicId)
}

func (s *favoriteService) addFavorite(userId int64, entityType string, entityId int64) errors.WTError {
	temp := s.GetBy(userId, entityType, entityId)
	if temp != nil { // 已经收藏
		return nil
	}
	return dao.FavoriteDao.Create(&model.Favorite{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: utils.NowTimestamp(),
	})
}
