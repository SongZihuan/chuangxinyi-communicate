package controller

import (
	"github.com/gin-gonic/gin"

	"gitee.com/wuntsong/chuangxinyi-communicate/oauth/gitee"
	"gitee.com/wuntsong/chuangxinyi-communicate/oauth/github"
	"gitee.com/wuntsong/chuangxinyi-communicate/oauth/qq"
)

// OAuthController oauth controller
type OAuthController struct {
	BaseController
}

// Authorize authorize
func (c *OAuthController) Authorize(ctx *gin.Context) {
	ref := ctx.Request.FormValue("ref")
	provider := ctx.Param("provider")
	params := map[string]string{"ref": ref}
	var url string
	if provider == "github" {
		url = github.AuthCodeURL(params)
	} else if provider == "gitee" {
		url = gitee.AuthCodeURL(params)
	} else {
		url = qq.AuthorizeUrl(params)
	}

	c.Success(ctx, gin.H{
		"url": url,
	})
}
