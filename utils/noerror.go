package utils

import (
	"github.com/pkg/errors"
)

type Logger interface {
	Error(string, ...any)
}

func MustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

func Recover(logger Logger, err *error, msg string) {
	e := recover()
	if e != nil {
		logger.Error("Error (%s): %v", msg, e)

		if err != nil {
			*err = errors.Errorf("respmsg (%s): %v", msg, e)
		}
	}
}
