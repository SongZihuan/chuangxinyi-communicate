package router

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/SongZihuan/chuangxinyi-communicate/controller"
	"github.com/SongZihuan/chuangxinyi-communicate/controller/admin"
	"github.com/SongZihuan/chuangxinyi-communicate/middleware"
)

var jwtAuth *jwt.GinJWTMiddleware

// Setup setup
func Setup(e *gin.Engine) {
	// 不使用gin的panic捕获，由access record捕获
	e.Use(Recovery())
	e.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Community API\n")
	})

	// e.Use(middleware.CurrentUser)

	//################################
	//#                              #
	//#             API              #
	//#                              #
	//################################
	api := e.Group("/api")

	api.Any("/stat", new(controller.SiteController).Stat)

	// JWT
	jwtAuth = middleware.JwtAuth(middleware.LoginStandard)
	api.POST("/auth/login", jwtAuth.LoginHandler)

	jwtApi := api.Group("/")
	jwtApi.Use(jwtAuth.MiddlewareFunc(), middleware.CurrentUser)

	// Configs
	configController := &controller.ConfigController{}
	api.GET("/config/configs", configController.List)

	// Nodes
	nodeController := &controller.NodeController{}
	api.GET("/nodes", nodeController.List)
	api.GET("/node/:id", nodeController.Show)

	// Topics
	topicController := &controller.TopicController{}
	api.GET("/topics", topicController.List)
	jwtApi.POST("/topics", topicController.Store)
	api.GET("/topic/:id", topicController.Show)

	jwtApi.GET("/topic/:id/edit", topicController.Edit)
	jwtApi.POST("/topic/:id/delete", topicController.Delete)
	jwtApi.PUT("/topic/:id", topicController.Update)

	api.GET("/topics/node", topicController.GetNodeTopics)
	api.GET("/topics/excellent", topicController.GetTopicsExcellent)
	api.GET("/topics/recommend", topicController.GetTopicsRecommend)
	api.GET("/topics/noreply", topicController.GetTopicsNoreply)
	api.GET("/topics/last", topicController.GetTopicsLast)
	api.GET("/topics/tag", topicController.GetTagTopics)
	api.GET("/topics/user/recent/:id", topicController.GetUserRecent)
	api.GET("/user/topics/:id", topicController.GetUserTopics)

	api.GET("/topic/:id/recentlikes", topicController.GetRecentLikes)
	jwtApi.POST("/topic/:id/like", topicController.Like)
	jwtApi.POST("/topic/:id/favorite", topicController.Favorite)

	// Comments
	commentController := &controller.CommentController{}
	api.GET("/comments", commentController.List)
	jwtApi.POST("/comments", commentController.Create)

	// Favorites
	favoriteController := &controller.FavoriteController{}
	jwtApi.GET("/favorites/favorited", favoriteController.GetFavorited)
	jwtApi.DELETE("/favorite/delete", favoriteController.Delete)

	// Tags
	tagController := &controller.TagController{}
	api.GET("/tag/:id", tagController.Show)
	api.GET("/tags", tagController.List)
	jwtApi.POST("/tag/autocomplete", tagController.Autocomplete)

	// Articles
	articleController := &controller.ArticleController{}
	jwtApi.POST("/articles", articleController.Store)
	api.GET("/article/:id", articleController.Show)
	jwtApi.GET("/article/:id/edit", articleController.Edit)
	jwtApi.POST("/article/:id/delete", articleController.Delete)
	jwtApi.PUT("/article/:id", articleController.Update)
	jwtApi.POST("/article/:id/favorite", articleController.Favorite)

	api.GET("/articles/user/newest/:id", articleController.GetUserNewestBy)

	api.GET("/articles/recent", articleController.GetRecent)
	api.GET("/articles/user/recent/:id", articleController.GetUserRecent)
	api.GET("/user/articles/:id", articleController.GetUserArticles)

	// header
	headerController := &controller.HeaderController{}
	api.GET("/header", headerController.GetHeader)

	// Users
	userController := &controller.UserController{}
	api.GET("/profile/:id", userController.Show)
	jwtApi.GET("/user/current", userController.GetCurrent)
	api.GET("/user/score/rank", userController.GetScoreRank)
	jwtApi.GET("/user/scorelogs", userController.GetScorelogs)
	jwtApi.GET("/user/favorites", userController.GetFavorites)

	api.GET("/users/:id/recentwatchers", userController.GetRecentWatchers)
	jwtApi.POST("/users/:id/watch", userController.Watch)
	jwtApi.GET("/watch/watched", userController.GetWatched)
	jwtApi.DELETE("/watch/delete", userController.WatchDelete)

	// Links
	linkController := &controller.LinkController{}
	api.GET("/links/top", linkController.GetToplinks)
	api.GET("/links", linkController.List)

	// Captcha
	captchaController := &controller.CaptchaController{}
	api.GET("/captcha/request", captchaController.GetRequest)
	api.GET("/captcha/show/:captchaId", captchaController.Show)

	// Upload/Get
	uploadController := &controller.UploadController{}
	jwtApi.POST("/upload", uploadController.Upload)
	jwtApi.POST("/upload/fetch", uploadController.UploadFromURL)
	api.GET("/file", uploadController.Get)

	//################################
	//#                              #
	//#          Admin API           #
	//#                              #
	//################################
	adminAPI := jwtApi.Group("/admin")
	adminAPI.Use(middleware.AdminRequired())

	// Dashboard
	dashboardController := &admin.DashboardController{}
	adminAPI.GET("/dashboard/systeminfo", dashboardController.Systeminfo)

	// Node
	adminNodeController := &admin.NodeController{}
	adminAPI.GET("/nodes", adminNodeController.List)
	adminAPI.GET("/nodes/:id", adminNodeController.Show)
	adminAPI.POST("/nodes", adminNodeController.Store)
	adminAPI.PUT("/nodes/:id", adminNodeController.Update)
	adminAPI.DELETE("/nodes/:id", adminNodeController.Delete)

	// Topic
	adminTopicController := &admin.TopicController{}
	adminAPI.GET("/topics", adminTopicController.List)
	adminAPI.GET("/topics/:id", adminTopicController.Show)
	adminAPI.PUT("/topics/:id", adminTopicController.Update)
	adminAPI.DELETE("/topics/:id", adminTopicController.Delete)
	adminAPI.POST("/topics/:id/recommend", adminTopicController.Recommend)
	adminAPI.POST("/topics/:id/unrecommend", adminTopicController.Unrecommend)
	adminAPI.POST("/topics/:id/undelete", adminTopicController.Undelete)

	// Article
	adminArticleController := &admin.ArticleController{}
	adminAPI.GET("/articles", adminArticleController.List)
	adminAPI.GET("/articles/:id", adminArticleController.Show)
	adminAPI.PUT("/articles/:id", adminArticleController.Update)
	adminAPI.DELETE("/articles/:id", adminArticleController.Delete)

	// Comment
	adminCommentController := &admin.CommentController{}
	adminAPI.GET("/comments", adminCommentController.List)
	adminAPI.GET("/comments/:id", adminCommentController.Show)
	adminAPI.PUT("/comments/:id", adminCommentController.Update)
	adminAPI.DELETE("/comments/:id", adminCommentController.Delete)

	// User
	adminUserController := &admin.UserController{}
	adminAPI.GET("/users", adminUserController.List)
	adminAPI.GET("/users/:id", adminUserController.Show)
	adminAPI.PUT("/users/:id", adminUserController.Update)
	adminAPI.DELETE("/users/:id", adminUserController.Delete)

	adminUserScoreController := &admin.UserScoreController{}
	adminAPI.GET("/user-scores", adminUserScoreController.List)
	adminAPI.GET("/user-scores/:id", adminUserScoreController.Show)

	adminUserScoreLogController := &admin.UserScoreLogController{}
	adminAPI.GET("/user-score-logs", adminUserScoreLogController.List)
	adminAPI.GET("/user-score-logs/:id", adminUserScoreLogController.Show)

	// Link
	adminLinkController := &admin.LinkController{}
	adminAPI.GET("/links", adminLinkController.List)
	adminAPI.GET("/links/:id", adminLinkController.Show)
	adminAPI.POST("/links", adminLinkController.Store)
	adminAPI.PUT("/links/:id", adminLinkController.Update)
	adminAPI.DELETE("/links/:id", adminLinkController.Delete)

	// Settings
	adminSettingController := &admin.SettingController{}
	adminAPI.GET("/settings", adminSettingController.List)
	adminAPI.POST("/settings", adminSettingController.Store)
}
