package init

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/redis"
)

func CloseAll() {
	_ = dao.Shutdown()
	redis.CloseRedis()
}
