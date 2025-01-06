package dao

import (
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

var RecordDao = newRecordDao()

func newRecordDao() *recordDao {
	return &recordDao{}
}

type recordDao struct {
}

func (d *recordDao) Get(id int64) *model.Record {
	ret := &model.Record{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *recordDao) Take(where ...interface{}) *model.Record {
	ret := &model.Record{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *recordDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Record) {
	cnd.Find(db, &list)
	return
}

func (d *recordDao) List(cnd *sqlcnd.SqlCnd) (list []model.Record, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Record{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *recordDao) Create(t *model.Record) (err errors.WTError) {
	err = errors.WarpQuick(db.Create(t).Error)
	return
}

func (d *recordDao) Update(t *model.Record) (err errors.WTError) {
	err = errors.WarpQuick(db.Save(t).Error)
	return
}

func (d *recordDao) Updates(id int64, columns map[string]interface{}) (err errors.WTError) {
	err = errors.WarpQuick(db.Model(&model.Record{}).Where("id = ?", id).Updates(columns).Error)
	return
}

func (d *recordDao) UpdateColumn(id int64, name string, value interface{}) (err errors.WTError) {
	err = errors.WarpQuick(db.Model(&model.Record{}).Where("id = ?", id).UpdateColumn(name, value).Error)
	return
}

func (d *recordDao) Delete(id int64) {
	db.Delete(&model.Record{}, "id = ?", id)
}

func (d *recordDao) DeleteOld(lastTime time.Time) (int64, errors.WTError) {
	res := db.Unscoped().Delete(&model.Record{}, "create_time < ?", lastTime)
	if res.Error != nil {
		return 0, errors.WarpQuick(res.Error)
	}
	return res.RowsAffected, nil
}
