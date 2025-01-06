package cache

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
)

type tagCache struct {
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	return &tagCache{}
}

func (c *tagCache) Get(tagId int64) *model.Tag {
	return dao.TagDao.Get(tagId)
}

func (c *tagCache) GetList(tagIds []int64) (tags []model.Tag) {
	if len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		tag := c.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (c *tagCache) Invalidate(tagId int64) {
	// 啥也不做
}
