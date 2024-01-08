package redis

import (
	"context"
	"gitee.com/wuntsong/chuangxinyi-communicate/signalexit"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

var client *redis.Client
var redcache *redis.Client
var memcache *cache.Cache

const KeepTTL = redis.KeepTTL

func InitRedis() errors.WTError {
	readdr := viper.GetString("redis.addr")
	reusername := viper.GetString("redis.userName")
	repassword := viper.GetString("redis.password")
	redb := viper.GetInt64("redis.db")

	client = redis.NewClient(&redis.Options{
		Addr:     readdr,
		Username: reusername,
		Password: repassword,
		DB:       int(redb),
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		_ = client.Close()
		client = nil
		return errors.Warp(err, "redis is fail to connect")
	}

	caaddr := viper.GetString("cache.addr")
	causername := viper.GetString("cache.userName")
	capassword := viper.GetString("cache.password")
	cadb := viper.GetInt64("cache.db")

	if len(caaddr) == 0 {
		redcache = nil
	} else {
		redcache = redis.NewClient(&redis.Options{
			Addr:     caaddr,
			Username: causername,
			Password: capassword,
			DB:       int(cadb),
		})

		err = redcache.Ping(context.Background()).Err()
		if err != nil {
			_ = client.Close()
			client = nil

			_ = redcache.Close()
			redcache = nil

			return errors.Warp(err, "cache is fail to connect")
		}
	}

	memcache = cache.New(time.Minute*30, time.Minute*60)

	signalexit.AddExitByFunc(lockExitFunc)

	return nil
}

func CloseRedis() {
	if client != nil {
		_ = client.Close()
		client = nil
	}

	if redcache != nil {
		_ = redcache.Close()
		redcache = nil
	}
}
