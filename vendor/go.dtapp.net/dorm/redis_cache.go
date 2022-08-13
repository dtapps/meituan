package dorm

import (
	"context"
	"encoding/json"
	"time"
)

// RedisCacheConfig 配置
type RedisCacheConfig struct {
	DefaultExpiration time.Duration // 过期时间
}

// RedisClientCache https://github.com/go-redis/redis
type RedisClientCache struct {
	config          *RedisCacheConfig
	operation       *RedisClient       // 操作
	GetterString    func() string      // 不存在的操作
	GetterInterface func() interface{} // 不存在的操作
}

// NewCache 实例化
func (c *RedisClient) NewCache(config *RedisCacheConfig) *RedisClientCache {
	cc := &RedisClientCache{config: config}
	cc.operation = c
	return cc
}

// NewCacheDefaultExpiration 实例化
func (c *RedisClient) NewCacheDefaultExpiration() *RedisClientCache {
	cc := &RedisClientCache{}
	cc.config.DefaultExpiration = time.Minute * 30 // 默认过期时间
	cc.operation = c
	return cc
}

// GetString 缓存操作
func (rc *RedisClientCache) GetString(ctx context.Context, key string) (ret string) {

	f := func() string {
		return rc.GetterString()
	}

	// 如果不存在，则调用GetterString
	ret, err := rc.operation.Get(ctx, key).Result()
	if err != nil {
		rc.operation.Set(ctx, key, f(), rc.config.DefaultExpiration)
		ret, _ = rc.operation.Get(ctx, key).Result()
	}

	return
}

// GetInterface 缓存操作
func (rc *RedisClientCache) GetInterface(ctx context.Context, key string, result interface{}) {

	f := func() string {
		marshal, _ := json.Marshal(rc.GetterInterface())
		return string(marshal)
	}

	// 如果不存在，则调用GetterInterface
	ret, err := rc.operation.Get(ctx, key).Result()

	if err != nil {
		rc.operation.Set(ctx, key, f(), rc.config.DefaultExpiration)
		ret, _ = rc.operation.Get(ctx, key).Result()
	}

	err = json.Unmarshal([]byte(ret), result)

	return
}
