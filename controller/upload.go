package controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SongZihuan/chuangxinyi-communicate/logger"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/SongZihuan/chuangxinyi-communicate/utils/uploader"
)

type UploadController struct {
	BaseController
}

const uploadMaxBytes int64 = 1024 * 1024 * 3 // 1M

// Upload upload file
func (c *UploadController) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	defer file.Close()
	if header.Size > uploadMaxBytes {
		c.Fail(ctx, utils.NewErrorMsg("图片不能超过3M"))
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}

	logger.Logger.Info("上传文件：%s, size: %s", header.Filename, header.Size)

	url, err := uploader.PutImage(fileBytes)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	data := make(map[string]string)
	data["url"] = url
	c.Success(ctx, data)
}

// Get get file
func (c *UploadController) Get(ctx *gin.Context) {
	key := ctx.Request.FormValue("key")
	if len(key) == 0 {
		c.Fail(ctx, utils.NewErrorMsg("key not give"))
		return
	}

	url, err := uploader.GetImage(key)
	if err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}
	c.Redirect(ctx, url)
}

// UploadFromURL fetch file by URL
func (c *UploadController) UploadFromURL(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, utils.ErrorNotLogin)
		return
	}

	data := make(map[string]string)
	if err := ctx.ShouldBindJSON(&data); err != nil {
		c.Fail(ctx, utils.FromError(err))
		return
	}

	url := data["url"]
	output, err := uploader.CopyImage(url)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"success": false,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}
	c.Success(ctx, gin.H{
		"originalURL": url,
		"url":         output,
	})
}
