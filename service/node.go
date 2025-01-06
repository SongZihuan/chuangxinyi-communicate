package service

import (
	errors "github.com/wuntsong-org/wterrors"
	"gorm.io/gorm"

	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

var NodeService = newNodeService()

func newNodeService() *nodeService {
	return &nodeService{}
}

type nodeService struct {
}

func (s *nodeService) Get(id int64) *model.Node {
	return dao.NodeDao.Get(id)
}

func (s *nodeService) List(cnd *sqlcnd.SqlCnd) (list []model.Node, paging *sqlcnd.Paging) {
	return dao.NodeDao.List(cnd)
}

func (s *nodeService) Create(dto form.NodeCreateForm) (*model.Node, errors.WTError) {
	node := &model.Node{
		Name:        dto.Name,
		Description: dto.Description,
		SortNo:      dto.SortNo,
		Status:      dto.Status,
		CreateTime:  utils.NowTimestamp(),
	}
	if err := dao.NodeDao.Create(node); err != nil {
		return nil, errors.New("创建节点失败")
	}
	cache.NodeCache.InvalidateAll()
	return node, nil
}

func (s *nodeService) Update(dto form.NodeUpdateForm) errors.WTError {
	err := dao.NodeDao.Updates(dto.ID, map[string]interface{}{
		"name":        dto.Name,
		"description": dto.Description,
		"sort_no":     dto.SortNo,
		"status":      dto.Status,
	})

	return errors.WarpQuick(err)
}

func (s *nodeService) Delete(id int64) {
	dao.NodeDao.Delete(id)
}

// 主题数+1
func (s *nodeService) IncrTopicCount(nodeId int64) {
	dao.DB().Model(&model.Node{}).Where("id = ?", nodeId).UpdateColumn("topic_count", gorm.Expr("topic_count + ?", 1))
}

func (s *nodeService) GetRecommendNodes() []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id").Limit(3))
}

func (s *nodeService) GetNodes() []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id"))
}
