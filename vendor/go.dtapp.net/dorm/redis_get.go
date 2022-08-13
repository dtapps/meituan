package dorm

import "github.com/go-redis/redis/v9"

// GetDb 获取驱动
func (c *RedisClient) GetDb() *redis.Client {
	return c.Db
}
