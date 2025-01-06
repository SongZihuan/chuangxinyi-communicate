package controller

import (
	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/convert"
	"github.com/SongZihuan/chuangxinyi-communicate/form"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/service"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/sqlcnd"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	BaseController
}

// GetCurrent get current user
func (c *UserController) GetCurrent(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"message": "ok",
		"data":    convert.ToSelfser(user),
	})
}

// 用户详情
func (c *UserController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		user := cache.UserCache.Get(gDto.ID)
		if user != nil && user.Status != model.StatusDeleted {
			c.Success(ctx, convert.ToUser(user))
		} else {
			c.Fail(ctx, utils.NewErrorMsg("用户不存在"))
		}
	}
}

// GetScoreRank 积分排行
func (c *UserController) GetScoreRank(ctx *gin.Context) {
	userScores := service.UserScoreService.Find(sqlcnd.NewSqlCnd().Desc("score").Limit(10))
	var results []*model.UserInfo
	for _, userScore := range userScores {
		results = append(results, convert.ToUserDefaultIfNull(userScore.UserId))
	}
	c.Success(ctx, results)
}

// GetScorelogs 用户积分记录
func (c *UserController) GetScorelogs(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	user := c.GetCurrentUser(ctx)

	logs, paging := service.UserScoreLogService.List(sqlcnd.NewSqlCnd().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	c.Success(ctx, gin.H{
		"results": logs,
		"paging":  paging,
	})
}

// GetFavorites get favorites
func (c *UserController) GetFavorites(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	cursor := form.FormValueInt64Default(ctx, "cursor", 0)

	// 查询列表
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = service.FavoriteService.Find(sqlcnd.NewSqlCnd().Where("user_id = ? and id < ?",
			user.ID, cursor).Desc("id").Limit(20))
	} else {
		favorites = service.FavoriteService.Find(sqlcnd.NewSqlCnd().Where("user_id = ?", user.ID).Desc("id").Limit(20))
	}

	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].ID
	}

	c.Success(ctx, gin.H{
		"results": convert.ToFavorites(favorites),
		"cursor":  cursor,
	})
}

// GetRecentWatchers 关注该用户的人
func (c *UserController) GetRecentWatchers(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		userWatchers := service.UserWatchService.Recent(gDto.ID, 10)
		var users []model.UserInfo
		for _, userWatcher := range userWatchers {
			userInfo := convert.ToUserById(userWatcher.WatcherID)
			if userInfo != nil {
				users = append(users, *userInfo)
			}
		}
		c.Success(ctx, users)
	}
}

// Watch 关注
func (c *UserController) Watch(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		err := service.UserWatchService.Watch(gDto.ID, user.ID)
		if err != nil {
			c.Fail(ctx, utils.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// GetWatched 是否关注了
func (c *UserController) GetWatched(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := form.FormValueInt64Default(ctx, "userId", 0)

	data := map[string]interface{}{}
	if user == nil || userID <= 0 {
		data["watched"] = false
	} else {
		tmp := service.UserWatchService.GetBy(userID, user.ID)
		data["watched"] = tmp != nil
	}
	c.Success(ctx, data)
}

// Delete 取消收藏
func (c *UserController) WatchDelete(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := form.FormValueInt64Default(ctx, "userId", 0)

	tmp := service.UserWatchService.GetBy(userID, user.ID)
	if tmp != nil {
		service.UserWatchService.Delete(tmp.ID)
	}
	c.Success(ctx, nil)
}
