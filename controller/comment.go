package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"gitee.com/wuntsong/chuangxinyi-communicate/convert"
	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

// CommentController comment controller
type CommentController struct {
	BaseController
}

// Create 发表评论
func (c *CommentController) Create(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var commentForm form.CommentCreateForm
	if c.BindAndValidate(ctx, &commentForm) {
		commentForm.UserID = user.ID
		comment, err := service.CommentService.Create(commentForm)
		if err != nil {
			c.Fail(ctx, utils.FromError(err))
			return
		}
		c.Success(ctx, convert.ToComment(*comment))
	}
}

// List 评论列表
func (c *CommentController) List(ctx *gin.Context) {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = form.FormValueInt64Default(ctx, "cursor", 0)

	entityType = ctx.Request.FormValue("entityType")
	if len(entityType) == 0 {
		c.Fail(ctx, &utils.CodeError{Message: "参数：entityType 不能为空"})
		return
	}

	if entityId, err = form.FormValueInt64(ctx, "entityId"); err != nil {
		c.Fail(ctx, &utils.CodeError{Message: err.Error()})
		return

	}

	comments, cursor := service.CommentService.GetComments(entityType, entityId, cursor)

	data := make(map[string]interface{})
	data["results"] = convert.ToComments(comments)
	data["cursor"] = strconv.FormatInt(cursor, 10)
	c.Success(ctx, data)
}
