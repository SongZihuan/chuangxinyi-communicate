package form

import (
	"fmt"
	"github.com/gin-gonic/gin"
	errors "github.com/wuntsong-org/wterrors"
	"strconv"
)

func FormValueInt(ctx *gin.Context, name string) (int, errors.WTError) {
	str := ctx.Request.FormValue(name)
	if str == "" {
		return 0, errors.New(fmt.Sprintf("unable to find param value '%s'", name))
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}
	return res, nil
}

func FormValueIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx *gin.Context, name string) (int64, errors.WTError) {
	str := ctx.Request.FormValue(name)
	if str == "" {
		return 0, errors.New(fmt.Sprintf("unable to find param value '%s'", name))
	}
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.WarpQuick(err)
	}
	return res, nil
}

func FormValueInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}
