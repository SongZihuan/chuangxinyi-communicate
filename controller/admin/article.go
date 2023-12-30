package admin

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/controller"
	"gitee.com/wuntsong/chuangxinyi-communicate/convert"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/html"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
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
	name := ctx.Request.FormValue("name")

	conditions := sqlcnd.NewSqlCnd()
	if len(name) > 0 {
		conditions.Like("name", name)
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
