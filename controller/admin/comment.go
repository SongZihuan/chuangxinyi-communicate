package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"gitee.com/wuntsong/chuangxinyi-communicate/controller"
	"gitee.com/wuntsong/chuangxinyi-communicate/convert"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/html"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
)

// CommentController comment controller
type CommentController struct {
	controller.BaseController
}

// Show show comment
func (c *CommentController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		comment := service.CommentService.Get(gDto.ID)
		if comment == nil {
			c.Fail(ctx, utils.NewErrorMsg("Comment not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, comment)
	}
}

// Update update a comment
func (c *CommentController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	comment := service.CommentService.Get(gDto.ID)
	if comment == nil {
		c.Fail(ctx, utils.NewErrorMsg("Comment not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var commentForm form.CommentUpdateForm
	if !c.BindAndValidate(ctx, &commentForm) {
		return
	}
	commentForm.ID = gDto.ID
	err := service.CommentService.Update(commentForm)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, comment)
}

// Delete delete comment
func (c *CommentController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.CommentService.Delete(gDto.ID)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// List list comments
func (c *CommentController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	status := ctx.Request.FormValue("status")
	entityType := ctx.Request.FormValue("entityType")
	entityID := ctx.Request.FormValue("entityId")
	userID := ctx.Request.FormValue("userId")

	conditions := sqlcnd.NewSqlCnd()
	if len(status) > 0 && status != "-1" {
		conditions.Eq("status", status)
	}
	if len(userID) > 0 {
		conditions.Eq("user_id", userID)
	}
	if len(entityType) > 0 {
		conditions.Eq("entity_type", entityType)
	}
	if len(entityID) > 0 {
		conditions.Eq("entity_id", entityID)
	}
	list, paging := service.CommentService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, comment := range list {
		result := utils.StructToMap(comment, "content")
		result["user"] = convert.ToUserDefaultIfNull(comment.UserId)
		mr := html.NewHtml().Run(comment.Content)
		result["content"] = mr.ContentHtml
		results = append(results, result)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}
