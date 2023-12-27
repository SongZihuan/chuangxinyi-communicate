package utils

import (
	"fmt"
	errors "github.com/wuntsong-org/wterrors"
	"strconv"
)

var (
	ErrorTopicNotFound    = NewError(-1, "话题不存在")
	ErrorArticleNotFound  = NewError(-2, "文章不存在")
	ErrorTagNotFound      = NewError(-3, "标签不存在")
	ErrorPermissionDenied = NewError(-100, "Permission denied.")
	ErrorNotLogin         = NewError(-200, "请先登录")
)

func NewError(code int, text string) *CodeError {
	return &CodeError{
		WTError: errors.New(text).SetCode(fmt.Sprintf("%d", code)),
	}
}

func NewErrorMsg(text string) *CodeError {
	return &CodeError{
		WTError: errors.New(text).SetCode("-1"),
	}
}

func NewErrorData(code int, text string, data interface{}) *CodeError {
	return &CodeError{
		WTError: errors.New(text).SetCode(fmt.Sprintf("%d", code)),
		Data:    data,
	}
}

func FromError(err error) *CodeError {
	if err == nil {
		return nil
	}

	var wterr errors.WTError
	if errors.As(err, &wterr) {
		return &CodeError{
			WTError: wterr,
		}
	}

	return &CodeError{
		WTError: errors.New(err.Error()).SetCode("-1"),
	}
}

type CodeError struct {
	errors.WTError
	Data interface{}
}

func (e *CodeError) CodeInt() int64 {
	code, err := strconv.ParseInt(e.Code(), 10, 64)
	if err != nil {
		return -1
	}
	return code
}
