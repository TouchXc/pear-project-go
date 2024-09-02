package dao

import (
	"context"
	"github.com/go-redis/redis"
	"ms_project/project-task/config"
	"time"
)

var RC *RedisCache
var IniConfig = config.InitConfig()
var redisOp = IniConfig.ReadRedisConfig()

type RedisCache struct {
	rdb *redis.Client
}

func (rc *RedisCache) Del(ctx context.Context, files []string) {
	rc.rdb.Del(files...)
}

func (rc *RedisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	result, err := rc.rdb.HKeys(key).Result()
	return result, err
}

func (rc *RedisCache) HSet(ctx context.Context, key string, field string, value string) {
	rc.rdb.HSet(key, field, value)
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(key, value, expire).Err()
	return err
}
func (rc *RedisCache) Get(ctx context.Context, key string) (value string, err error) {
	value, err = rc.rdb.Get(key).Result()
	return
}
func InitRedis() {
	rdb := redis.NewClient(redisOp)
	RC = &RedisCache{
		rdb: rdb,
	}
}
