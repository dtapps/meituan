package dorm

import (
	"context"
	"github.com/go-redis/redis/v9"
)

// Subscribe 订阅channel
func (c *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.Db.Subscribe(ctx, channels...)
}

// PSubscribe 订阅channel支持通配符匹配
func (c *RedisClient) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.Db.PSubscribe(ctx, channels...)
}

// Publish 将信息发送到指定的channel
func (c *RedisClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.Db.Publish(ctx, channel, message)
}

// PubSubChannels 查询活跃的channel
func (c *RedisClient) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.Db.PubSubChannels(ctx, pattern)
}

// PubSubNumSub 查询指定的channel有多少个订阅者
func (c *RedisClient) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	return c.Db.PubSubNumSub(ctx, channels...)
}
