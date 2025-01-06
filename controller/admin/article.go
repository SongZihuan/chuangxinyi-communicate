package admin

import (
	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/controller"
	"github.com/SongZihuan/chuangxinyi-communicate/convert"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/html"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ArticleController article controller
type ArticleController struct {
	controller.BaseController
}

// Show show article
func (c *ArticleController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		article := service.ArticleService.Get(gDto.ID)
		if article == nil {
			c.Fail(ctx, utils.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, article)
	}
}

// Update update a article
func (c *ArticleController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	article := service.ArticleService.Get(gDto.ID)
	if article == nil {
		c.Fail(ctx, utils.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var articleForm form.ArticleUpdateForm
	if !c.BindAndValidate(ctx, &articleForm) {
		return
	}
	articleForm.ID = gDto.ID
	err := service.ArticleService.Update(articleForm)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, article)
}

// Delete delete article
func (c *ArticleController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.ArticleService.Delete(gDto.ID)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// List list articles
func (c *ArticleController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	userID := ctx.Request.FormValue("user_id")
	status := ctx.Request.FormValue("status")
	title := ctx.Request.FormValue("title")

	conditions := sqlcnd.NewSqlCnd()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(userID) > 0 {
		conditions.Eq("user_id", userID)
	}
	if len(status) > 0 && status != "-1" {
		conditions.Eq("status", status)
	}
	if len(title) > 0 {
		conditions.Like("title", title)
	}

	list, paging := service.ArticleService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, article := range list {
		item := utils.StructToMap(article, "content")
		item["user"] = convert.ToUserDefaultIfNull(article.UserId)

		// 简介
		mr := html.NewHtml().Run(article.Content)
		if len(article.Summary) == 0 {
			item["summary"] = mr.SummaryText
		}

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.ID)
		tags := cache.TagCache.GetList(tagIds)
		item["tags"] = convert.ToTags(tags)

		results = append(results, item)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}
