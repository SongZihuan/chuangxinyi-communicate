package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/SongZihuan/chuangxinyi-communicate/controller"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

// LinkController link controller
type LinkController struct {
	controller.BaseController
}

// Show show link
func (c *LinkController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		link := service.LinkService.Get(gDto.ID)
		if link == nil {
			c.Fail(ctx, utils.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, link)
	}
}

// Store create a link
func (c *LinkController) Store(ctx *gin.Context) {
	var linkForm form.LinkCreateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	link, err := service.LinkService.Create(linkForm)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Update update a link
func (c *LinkController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	link := service.LinkService.Get(gDto.ID)
	if link == nil {
		c.Fail(ctx, utils.NewErrorMsg("Link not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var linkForm form.LinkUpdateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	linkForm.ID = gDto.ID
	err := service.LinkService.Update(linkForm)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Delete delete link
func (c *LinkController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.LinkService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list links
func (c *LinkController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)

	title := ctx.Request.FormValue("title")
	url := ctx.Request.FormValue("url")
	status := ctx.Request.FormValue("status")

	conditions := sqlcnd.NewSqlCnd()
	if len(status) > 0 && status != "-1" {
		conditions.Eq("status", status)
	}
	if len(url) > 0 {
		conditions.Eq("url", url)
	}
	if len(title) > 0 {
		conditions.Eq("title", title)
	}
	list, paging := service.LinkService.List(conditions.Page(page, limit).Desc("id"))

	c.Success(ctx, &sqlcnd.PageResult{Results: list, Page: paging})
}
