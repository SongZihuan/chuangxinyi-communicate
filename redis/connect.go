package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

var RedisClient *redis.Client

func InitRedis() errors.WTError {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	db := viper.GetInt64("redis.db")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(db),
	})

	err := RedisClient.Ping(context.Background()).Err()
	if err != nil {
		_ = RedisClient.Close()
		RedisClient = nil
		return errors.Errorf("redis is fail to connect: %s", err.Error())
	}

	return nil
}

func CloseRedis() {
	if RedisClient == nil {
		return
	}

	_ = RedisClient.Close()
	RedisClient = nil
}
