package init

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/redis"
)

func CloseAll() {
	_ = dao.Shutdown()
	redis.CloseRedis()
}
