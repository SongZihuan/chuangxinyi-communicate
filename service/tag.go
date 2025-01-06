package service

import (
	"fmt"
	"github.com/SongZihuan/chuangxinyi-communicate/yundun"
	errors "github.com/wuntsong-org/wterrors"
	"strings"

	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

type ScanTagCallback func(tags []model.Tag) bool

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (s *tagService) Get(id int64) *model.Tag {
	return dao.TagDao.Get(id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return dao.TagDao.Take(where...)
}

func (s *tagService) Find(cnd *sqlcnd.SqlCnd) []model.Tag {
	return dao.TagDao.Find(cnd)
}

func (s *tagService) FindOne(cnd *sqlcnd.SqlCnd) *model.Tag {
	return dao.TagDao.FindOne(cnd)
}

func (s *tagService) List(cnd *sqlcnd.SqlCnd) (list []model.Tag, paging *sqlcnd.Paging) {
	return dao.TagDao.List(cnd)
}

func (s *tagService) Create(t *model.Tag) errors.WTError {
	ok, err := yundun.CheckText(fmt.Sprintf("标签：%s", t.Name))
	if err != nil {
		return errors.WarpQuick(err)
	} else if !ok {
		return errors.Errorf("内容不合法")
	}

	return dao.TagDao.Create(t)
}

func (s *tagService) Update(t *model.Tag) errors.WTError {
	ok, err := yundun.CheckText(fmt.Sprintf("标签：%s", t.Name))
	if err != nil {
		return errors.WarpQuick(err)
	} else if !ok {
		return errors.Errorf("内容不合法")
	}

	if err := dao.TagDao.Update(t); err != nil {
		return errors.WarpQuick(err)
	}
	cache.TagCache.Invalidate(t.ID)
	return nil
}

// 自动完成
func (s *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("status = ? and name like ?",
		model.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, errors.WTError) {
	return dao.TagDao.GetOrCreate(name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return dao.TagDao.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("status = ?", model.StatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagId: tag.ID, TagName: tag.Name})
	}
	return tags
}

func (s *tagService) GetTagInIds(tagIds []int64) []model.Tag {
	return dao.TagDao.GetTagInIds(tagIds)
}

// 扫描
func (s *tagService) Scan(cb ScanTagCallback) {
	var cursor int64
	for {
		list := dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		if !cb(list) {
			break
		}
	}
}
