package core

import (
	"github.com/go-redis/redis/v8"
	"petHealthToolApi/config"
	"petHealthToolApi/global"
)

// redis初始化

var (
	RedisDb *redis.Client
)

func init() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Address,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.Db,
	})
	_, err := RedisDb.Ping(global.Ctx).Result()
	if err != nil {
		return
	}
	global.Log.Infof("Redis init success")
	return
}

// GetRedisDb 获取redis实例
func GetRedisDb() *redis.Client {
	if RedisDb == nil {
		panic("redis init failed")
	}
	return RedisDb
}
