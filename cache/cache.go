package cache

import (
	"github.com/goburrow/cache"
	errors "github.com/wuntsong-org/wterrors"
)

func key2Int64(key cache.Key) int64 {
	return key.(int64)
}

func Setup() errors.WTError {
	return nil
}
