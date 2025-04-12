package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type config struct {
		Addr string
	}
	var conf config
	err := viper.UnmarshalKey("cache.redis", &conf)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return redisClient
}
