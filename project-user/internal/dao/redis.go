package dao

import (
	"context"
	"github.com/go-redis/redis"
	"ms_project/project-user/config"
	"time"
)

var RC *RedisCache
var IniConfig = config.InitConfig()
var redisOp = IniConfig.ReadRedisConfig()

type RedisCache struct {
	rdb *redis.Client
}

func InitRedis() {
	rdb := redis.NewClient(redisOp)
	RC = &RedisCache{
		rdb: rdb,
	}
}
func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(key, value, expire).Err()
	return err
}
func (rc *RedisCache) Get(ctx context.Context, key string) (value string, err error) {
	value, err = rc.rdb.Get(key).Result()
	return
}
