package cache

import (
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/zedisdog/armor/config"
)

//RedClient 连接引用
var RedClient *redis.Client

//Instance 获取实例
func Instance() *redis.Client {
	return RedClient
}

//InitCache 初始化缓存连接
func InitCache() (*redis.Client, error) {
	//log.Log.Infoln(config.Instance().String("cache.redis.host") + ":" + config.Instance().String("cache.redis.port"))
	RedClient = redis.NewClient(&redis.Options{
		Addr:     config.Instance().String("cache.redis.host") + ":" + strconv.Itoa(config.Instance().Int("cache.redis.port")),
		Network:  config.Instance().String("cache.redis.network"),
		Password: config.Instance().String("cache.redis.password"),
		DB:       config.Instance().Int("cache.redis.db"),
		PoolSize: config.Instance().Int("cache.redis.poolSize"),
	})
	_, err := RedClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return RedClient, nil
}
