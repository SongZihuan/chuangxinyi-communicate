package controller

import (
	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/convert"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

type NodeController struct {
	BaseController
}

// List 节点列表
func (c *NodeController) List(ctx *gin.Context) {
	nodes := cache.NodeCache.GetAll()

	c.Success(ctx, convert.ToNodes(nodes))
}

// Show 显示单个节点
func (c *NodeController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		node := service.NodeService.Get(gDto.ID)
		c.Success(ctx, convert.ToNode(node))
	}
}
