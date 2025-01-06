package service

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
	errors "github.com/wuntsong-org/wterrors"
)

var UserScoreLogService = newUserScoreLogService()

func newUserScoreLogService() *userScoreLogService {
	return &userScoreLogService{}
}

type userScoreLogService struct {
}

func (s *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return dao.UserScoreLogDao.Get(id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return dao.UserScoreLogDao.Take(where...)
}

func (s *userScoreLogService) Find(cnd *sqlcnd.SqlCnd) []model.UserScoreLog {
	return dao.UserScoreLogDao.Find(cnd)
}

func (s *userScoreLogService) FindOne(cnd *sqlcnd.SqlCnd) *model.UserScoreLog {
	return dao.UserScoreLogDao.FindOne(cnd)
}

func (s *userScoreLogService) List(cnd *sqlcnd.SqlCnd) (list []model.UserScoreLog, paging *sqlcnd.Paging) {
	return dao.UserScoreLogDao.List(cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) errors.WTError {
	return dao.UserScoreLogDao.Create(t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) errors.WTError {
	return dao.UserScoreLogDao.Update(t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) errors.WTError {
	return dao.UserScoreLogDao.Updates(id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) errors.WTError {
	return dao.UserScoreLogDao.UpdateColumn(id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	dao.UserScoreLogDao.Delete(id)
}
