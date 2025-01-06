package service

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
	errors "github.com/wuntsong-org/wterrors"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (s *topicTagService) Get(id int64) *model.TopicTag {
	return dao.TopicTagDao.Get(id)
}

func (s *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return dao.TopicTagDao.Take(where...)
}

func (s *topicTagService) Find(cnd *sqlcnd.SqlCnd) []model.TopicTag {
	return dao.TopicTagDao.Find(cnd)
}

func (s *topicTagService) FindOne(cnd *sqlcnd.SqlCnd) *model.TopicTag {
	return dao.TopicTagDao.FindOne(cnd)
}

func (s *topicTagService) List(cnd *sqlcnd.SqlCnd) (list []model.TopicTag, paging *sqlcnd.Paging) {
	return dao.TopicTagDao.List(cnd)
}

func (s *topicTagService) Create(t *model.TopicTag) errors.WTError {
	return dao.TopicTagDao.Create(t)
}

func (s *topicTagService) Update(t *model.TopicTag) errors.WTError {
	return dao.TopicTagDao.Update(t)
}

func (s *topicTagService) Updates(id int64, columns map[string]interface{}) errors.WTError {
	return dao.TopicTagDao.Updates(id, columns)
}

func (s *topicTagService) UpdateColumn(id int64, name string, value interface{}) errors.WTError {
	return dao.TopicTagDao.UpdateColumn(id, name, value)
}

func (s *topicTagService) DeleteByTopicId(topicId int64) {
	dao.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", model.StatusDeleted)
}

func (s *topicTagService) UndeleteByTopicId(topicId int64) {
	dao.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", model.StatusOk)
}
