package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/convert"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

type TagController struct {
	BaseController
}

// Show 标签详情
func (c *TagController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		tag := cache.TagCache.Get(gDto.ID)
		if tag == nil {
			c.Fail(ctx, utils.ErrorTagNotFound)
			return
		}
		c.Success(ctx, convert.ToTag(tag))
	}
}

// List 标签列表
func (c *TagController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	tags, paging := service.TagService.List(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).
		Page(page, 200).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = convert.ToTags(tags)
	data["page"] = paging
	c.Success(ctx, data)
}

// Autocomplete 标签自动完成
func (c *TagController) Autocomplete(ctx *gin.Context) {
	input := ctx.Request.FormValue("input")
	tags := service.TagService.Autocomplete(input)
	c.Success(ctx, tags)
}
