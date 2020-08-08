package cache

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"strconv"

	"github.com/go-redis/redis/v7"
)

func New(v *viper.Viper) (*redis.Client, error) {
	if v.GetBool("cache.redis.enable") {
		RedClient := redis.NewClient(&redis.Options{
			Addr:     v.GetString("cache.redis.host") + ":" + strconv.Itoa(v.GetInt("cache.redis.port")),
			Network:  v.GetString("cache.redis.network"),
			Password: v.GetString("cache.redis.password"),
			DB:       v.GetInt("cache.redis.db"),
			PoolSize: v.GetInt("cache.redis.poolSize"),
		})
		if _, err := RedClient.Ping().Result(); err != nil {
			return nil, err
		}

		return RedClient, nil
	}

	return nil, nil
}

var ProviderSet = wire.NewSet(New)
