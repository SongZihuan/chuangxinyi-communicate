package rand

import (
	mathrand "math/rand"
	"time"
)

var GlobalRander *mathrand.Rand

func InitRander() error {
	GlobalRander = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	return nil
}
