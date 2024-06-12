package meituan

import (
	"go.dtapp.net/gorequest"
)

// ConfigClient 配置
func (c *Client) ConfigClient(config *ClientConfig) {
	c.config.secret = config.Secret
	c.config.appKey = config.AppKey
}

func (c *Client) SetSecret(secret string) *Client {
	c.config.secret = secret
	return c
}

func (c *Client) SetAppKey(appKey string) *Client {
	c.config.appKey = appKey
	return c
}

// SetClientIP 配置
func (c *Client) SetClientIP(clientIP string) *Client {
	c.config.clientIP = clientIP
	if c.httpClient != nil {
		c.httpClient.SetClientIP(clientIP)
	}
	return c
}

// SetLogFun 设置日志记录函数
func (c *Client) SetLogFun(logFun gorequest.LogFunc) {
	if c.httpClient != nil {
		c.httpClient.SetLogFunc(logFun)
	}
}
