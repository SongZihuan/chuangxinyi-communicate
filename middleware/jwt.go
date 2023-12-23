package middleware

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/auth/login"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"gitee.com/wuntsong/chuangxinyi-communicate/form"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/log"
)

// login type
var (
	LoginStandard = 1
	LoginOAuth    = 2
)

type LoginDto struct {
	LoginToken string `form:"loginToken" json:"loginToken" binding:"required"`
}

// LoginOAuthDto - oauth login
type LoginOAuthDto struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// todo : 用单独的claims model去掉user model
func JwtAuth(LoginType int) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "Jwt",
		// SigningAlgorithm: "RS256",
		// PubKeyFile:       "keys/jwt_private_key.pem",
		// PrivKeyFile:      "keys/jwt_public_key.pem",
		Key:           []byte(viper.GetString("jwt.key")),
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour * 24 * 90,
		IdentityKey:   viper.GetString("jwt.identity_key"),
		LoginResponse: LoginResponse,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(model.UserClaims); ok {
				return jwt.MapClaims{
					"id":    v.ID,
					"name":  utils.GetUserName(v.Phone, v.Email, v.Username, v.Nickname),
					"uid":   v.Uid,
					"uname": v.Phone,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			user := dao.UserDao.GetByUid(claims["uid"].(string))
			if user == nil {
				return nil
			}

			return model.UserClaims{
				User: user,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return Authenticator(c)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(model.UserClaims); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(200, gin.H{
				"code":    200,
				"success": false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return jwtMiddleware
}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": map[string]interface{}{
			"token":  token,
			"expire": expire,
		},
		"success": true,
		"message": "success",
	})
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginDto LoginDto
	if err := form.Bind(c, &loginDto); err != nil {
		return "", err
	}

	user, err := login.CheckLogin(c, loginDto.LoginToken)
	if err != nil {
		return nil, err
	}

	return model.UserClaims{
		User: user,
	}, nil
}
