package middleware

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/api/auth/login" {
			logger.Logger.Tag("A", c.Request.Method)
		}
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)      // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Credentials", "true") // 可将将 * 替换为指定的域名
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}
