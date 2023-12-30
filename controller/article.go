package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"gitee.com/wuntsong/chuangxinyi-communicate/convert"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
)

type ArticleController struct {
	BaseController
}

// Show show article by id
func (c *ArticleController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		article := service.ArticleService.Get(gDto.ID)
		if article == nil || article.Status != model.StatusOk {
			c.Fail(ctx, utils.ErrorArticleNotFound)
			return
		}
		c.Success(ctx, convert.ToArticle(article))
	}
}

// Create 发表文章
func (c *ArticleController) Store(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, utils.ErrorNotLogin)
		return
	}
	var articleForm form.ArticleCreateForm
	if c.BindAndValidate(ctx, &articleForm) {
		articleForm.UserID = user.ID
		article, err := service.ArticleService.Create(articleForm)
		if err != nil {
			c.Fail(ctx, utils.FromError(err))
			return
		}
		c.Success(ctx, convert.ToArticle(article))
	}
}

// Edit 编辑时获取详情
func (c *ArticleController) Edit(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, utils.ErrorNotLogin)
		return
	}

	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		article := service.ArticleService.Get(gDto.ID)

		if article == nil || article.Status != model.StatusOk {
			c.Fail(ctx, utils.NewErrorMsg("话题不存在或已被删除"))
			return
		}
		if article.UserId != user.ID {
			c.Fail(ctx, utils.NewErrorMsg("无权限"))
			return
		}

		tags := service.ArticleService.GetArticleTags(article.ID)
		var tagNames []string
		if len(tags) > 0 {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		c.Success(ctx, gin.H{
			"articleId": article.ID,
			"title":     article.Title,
			"content":   article.Content,
			"tags":      tagNames,
		})
	}
}

// Update 编辑文章
func (c *ArticleController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, utils.ErrorNotLogin)
		return
	}
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, utils.ErrorArticleNotFound)
		return
	}

	article := service.ArticleService.Get(gDto.ID)
	if article == nil || article.Status == model.StatusDeleted {
		c.Fail(ctx, utils.ErrorArticleNotFound)
		return
	}

	if article.UserId != user.ID {
		c.Fail(ctx, utils.NewErrorMsg("无权限"))
		return
	}

	var articleForm form.ArticleUpdateForm
	if c.BindAndValidate(ctx, &articleForm) {
		articleForm.ID = article.ID
		err := service.ArticleService.Update(articleForm)
		if err != nil {
			c.Fail(ctx, utils.FromError(err))
			return
		}
		c.Success(ctx, gin.H{
			"articleId": article.ID,
		})
	}
}

// Delete 删除文章
func (c *ArticleController) Delete(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, utils.ErrorNotLogin)
		return
	}
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, utils.ErrorArticleNotFound)
		return
	}

	article := service.ArticleService.Get(gDto.ID)
	if article == nil || article.Status == model.StatusDeleted {
		c.Fail(ctx, utils.ErrorArticleNotFound)
		return
	}

	if article.UserId != user.ID {
		c.Fail(ctx, utils.NewErrorMsg("无权限"))
		return
	}

	err := service.ArticleService.Delete(article.ID)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// GetRecent 最近文章
func (c *ArticleController) GetRecent(ctx *gin.Context) {
	articles := service.ArticleService.Find(sqlcnd.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(10))

	c.Success(ctx, articles)
}

// List 文章列表
func (c *ArticleController) List(ctx *gin.Context) {
	cursor := form.FormValueInt64Default(ctx, "cursor", 0)
	articles, cursor := service.ArticleService.GetArticles(cursor)
	c.Success(ctx, gin.H{
		"results": convert.ToSimpleArticles(articles),
		"cursor":  strconv.FormatInt(cursor, 10),
	})
}

// GetTagArticles 标签文章列表
func (c *ArticleController) GetTagArticles(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		cursor := form.FormValueInt64Default(ctx, "cursor", 0)
		articles, cursor := service.ArticleService.GetTagArticles(gDto.ID, cursor)
		c.Success(ctx, gin.H{
			"results": convert.ToSimpleArticles(articles),
			"cusor":   strconv.FormatInt(cursor, 10),
		})
	}
}

// GetUserRecent 用户最近的文章
func (c *ArticleController) GetUserRecent(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		articles := service.ArticleService.Find(sqlcnd.NewSqlCnd().Where("user_id = ? and status = ?",
			gDto.ID, model.StatusOk).Desc("id").Limit(10))
		c.Success(ctx, convert.ToSimpleArticles(articles))
	}
}

// GetUserArticles 用户的文章
func (c *ArticleController) GetUserArticles(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		articles, paging := service.ArticleService.List(sqlcnd.NewSqlCnd().
			Eq("user_id", gDto.ID).
			Eq("status", model.StatusOk).
			Page(page, 20).Desc("id"))

		c.Success(ctx, gin.H{
			"results": convert.ToSimpleArticles(articles),
			"page":    paging,
		})
	}
}

// GetUserNewestBy 用户最新的文章
func (c *ArticleController) GetUserNewestBy(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		newestArticles := service.ArticleService.GetUserNewestArticles(gDto.ID)
		c.Success(ctx, convert.ToSimpleArticles(newestArticles))
	}
}

// GetRelatedBy 相关文章
func (c *ArticleController) GetRelatedBy(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		relatedArticles := service.ArticleService.GetRelatedArticles(gDto.ID)
		c.Success(ctx, convert.ToSimpleArticles(relatedArticles))
	}
}

// Favorite 收藏文章
func (c *ArticleController) Favorite(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, utils.ErrorArticleNotFound)
		return
	}
	err := service.FavoriteService.AddArticleFavorite(user.ID, gDto.ID)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
