package meituan

import (
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel"
)

// ClientConfig 实例配置
type ClientConfig struct {
	ClientIP string // 客户端IP
	Secret   string // 秘钥
	AppKey   string // 渠道标记
}

// Client 实例
type Client struct {
	httpClient *gorequest.App
	config     struct {
		clientIP string // 客户端IP
		secret   string // 秘钥
		appKey   string // 渠道标记
	}
}

// NewClient 创建实例化
func NewClient(config *ClientConfig) (*Client, error) {

	c := &Client{}

	c.httpClient = gorequest.NewHttp()
	c.httpClient.SetTracer(otel.Tracer("go.dtapp.net/meituan"))

	c.config.clientIP = config.ClientIP
	c.config.secret = config.Secret
	c.config.appKey = config.AppKey

	return c, nil
}
