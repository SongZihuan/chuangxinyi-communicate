package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/SongZihuan/chuangxinyi-communicate/controller"
	"github.com/SongZihuan/chuangxinyi-communicate/convert"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
)

// UserScoreLogController user score controller
type UserScoreLogController struct {
	controller.BaseController
}

// Show 显示积分纪录
func (c *UserScoreLogController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		userScoreLog := service.UserScoreLogService.Get(gDto.ID)
		if userScoreLog == nil {
			c.Fail(ctx, utils.NewErrorMsg("User score log not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, userScoreLog)
	}
}

// List 显示积分列表
func (c *UserScoreLogController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	userId := ctx.Request.FormValue("userId")
	sourceType := ctx.Request.FormValue("sourceType")
	sourceId := ctx.Request.FormValue("sourceId")
	ltype := ctx.Request.FormValue("type")

	conditions := sqlcnd.NewSqlCnd()
	if len(userId) > 0 {
		conditions.Eq("user_id", userId)
	}
	if len(sourceType) > 0 {
		conditions.Eq("source_type", sourceType)
	}
	if len(sourceId) > 0 {
		conditions.Eq("source_id", sourceId)
	}
	if len(ltype) > 0 {
		conditions.Eq("type", ltype)
	}

	list, paging := service.UserScoreLogService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		item := utils.StructToMap(userScoreLog)
		item["user"] = convert.ToUserDefaultIfNull(userScoreLog.UserId)
		results = append(results, item)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}
