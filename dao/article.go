package dao

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
	errors "github.com/wuntsong-org/wterrors"
)

var ArticleDao = newArticleDao()

func newArticleDao() *articleDao {
	return &articleDao{}
}

type articleDao struct {
}

// Get returns article by given ID.
func (d *articleDao) Get(id int64) *model.Article {
	ret := &model.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Article) {
	cnd.Find(db, &list)
	return
}

func (d *articleDao) List(cnd *sqlcnd.SqlCnd) (list []model.Article, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *articleDao) Create(t *model.Article) (err errors.WTError) {
	err = errors.WarpQuick(db.Create(t).Error)
	return
}

func (d *articleDao) Update(t *model.Article) (err errors.WTError) {
	err = errors.WarpQuick(db.Save(t).Error)
	return
}

func (d *articleDao) Updates(id int64, columns map[string]interface{}) (err errors.WTError) {
	err = errors.WarpQuick(db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error)
	return
}

func (d *articleDao) UpdateColumn(id int64, name string, value interface{}) (err errors.WTError) {
	err = errors.WarpQuick(db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error)
	return
}

func (d *articleDao) Delete(id int64) {
	db.Delete(&model.Article{}, "id = ?", id)
}
