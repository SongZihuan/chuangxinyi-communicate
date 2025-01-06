package service

import (
	errors "github.com/wuntsong-org/wterrors"
	"gorm.io/gorm"

	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (s *linkService) Get(id int64) *model.Link {
	return dao.LinkDao.Get(id)
}

func (s *linkService) Find(cnd *sqlcnd.SqlCnd) []model.Link {
	return dao.LinkDao.Find(cnd)
}

func (s *linkService) List(cnd *sqlcnd.SqlCnd) (list []model.Link, paging *sqlcnd.Paging) {
	return dao.LinkDao.List(cnd)
}

func (s *linkService) Create(dto form.LinkCreateForm) (*model.Link, errors.WTError) {
	link := &model.Link{
		Title:      dto.Title,
		Url:        dto.URL,
		Summary:    dto.Summary,
		Logo:       dto.Logo,
		CreateTime: utils.NowTimestamp(),
	}
	err := dao.Tx(dao.DB(), func(tx *gorm.DB) errors.WTError {
		err := dao.LinkDao.Create(link)
		if err != nil {
			return errors.WarpQuick(err)
		}
		return nil
	})
	return link, errors.WarpQuick(err)
}

func (s *linkService) Update(dto form.LinkUpdateForm) errors.WTError {
	err := dao.LinkDao.Updates(dto.ID, map[string]interface{}{
		"title":       dto.Title,
		"url":         dto.URL,
		"summary":     dto.Summary,
		"logo":        dto.Logo,
		"status":      dto.Status,
		"update_time": utils.NowTimestamp(),
	})

	return errors.WarpQuick(err)
}

func (s *linkService) Delete(id int64) {
	dao.LinkDao.Delete(id)
}
