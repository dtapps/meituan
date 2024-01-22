package meituan

import (
	"go.dtapp.net/golog"
)

// ClientConfig 实例配置
type ClientConfig struct {
	Secret string // 秘钥
	AppKey string // 渠道标记
}

// Client 实例
type Client struct {
	config struct {
		secret string // 秘钥
		appKey string // 渠道标记
	}
	gormLog struct {
		status bool           // 状态
		client *golog.ApiGorm // 日志服务
	}
	mongoLog struct {
		status bool            // 状态
		client *golog.ApiMongo // 日志服务
	}
}

// NewClient 创建实例化
func NewClient(config *ClientConfig) (*Client, error) {

	c := &Client{}

	c.config.secret = config.Secret
	c.config.appKey = config.AppKey

	return c, nil
}
