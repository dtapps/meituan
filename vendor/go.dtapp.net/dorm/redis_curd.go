package dorm

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

// Set 设置一个key的值
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.Db.Set(ctx, key, value, expiration)
}

// Get 查询key的值
func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.Db.Get(ctx, key)
}

// GetSet 设置一个key的值，并返回这个key的旧值
func (c *RedisClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	return c.Db.GetSet(ctx, key, value)
}

// SetNX 如果key不存在，则设置这个key的值
func (c *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.Db.SetNX(ctx, key, value, expiration)
}

// MGet 批量查询key的值
func (c *RedisClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return c.Db.MGet(ctx, keys...)
}

// MSet 批量设置key的值
// MSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
func (c *RedisClient) MSet(ctx context.Context, values map[string]interface{}) *redis.StatusCmd {
	return c.Db.MSet(ctx, values)
}

// Incr 针对一个key的数值进行递增操作
func (c *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.Db.Incr(ctx, key)
}

// IncrBy 针对一个key的数值进行递增操作，指定每次递增多少
func (c *RedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return c.Db.IncrBy(ctx, key, value)
}

// Decr 针对一个key的数值进行递减操作
func (c *RedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	return c.Db.Decr(ctx, key)
}

// DecrBy 针对一个key的数值进行递减操作，指定每次递减多少
func (c *RedisClient) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return c.Db.DecrBy(ctx, key, value)
}

// Del 删除key操作，支持批量删除
func (c *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.Db.Del(ctx, keys...)
}

// Keys 按前缀获取所有 key
func (c *RedisClient) Keys(ctx context.Context, prefix string) *redis.SliceCmd {
	values, _ := c.Db.Keys(ctx, prefix).Result()
	if len(values) <= 0 {
		return &redis.SliceCmd{}
	}
	keys := make([]string, 0, len(values))
	for _, value := range values {
		keys = append(keys, value)
	}
	return c.MGet(ctx, keys...)
}
