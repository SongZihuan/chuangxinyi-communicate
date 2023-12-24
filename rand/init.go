package rand

import (
	errors "github.com/wuntsong-org/wterrors"
	mathrand "math/rand"
	"time"
)

var GlobalRander *mathrand.Rand

func InitRander() errors.WTError {
	GlobalRander = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	return nil
}
