package init

import "gitee.com/wuntsong/chuangxinyi-communicate/redis"

func CloseAll() {
	redis.CloseRedis()
}
