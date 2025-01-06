package cache

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

var (
	allNodesCacheKey = "all_nodes_cache"
)

type nodeCache struct {
}

var NodeCache = newNodeCache()

func newNodeCache() *nodeCache {
	return &nodeCache{}
}

func (c *nodeCache) Get(nodeId int64) *model.Node {
	if nodeId <= 0 {
		return nil
	}
	return dao.NodeDao.Get(nodeId)
}

func (c *nodeCache) Invalidate(nodeId int64) {
	// 什么都不做
}

func (c *nodeCache) GetAll() []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id"))
}

func (c *nodeCache) InvalidateAll() {
	// 什么都不做
}
