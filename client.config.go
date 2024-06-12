package meituan

import (
	"go.dtapp.net/gorequest"
)

// ConfigClient 配置
func (c *Client) ConfigClient(config *ClientConfig) {
	c.config.secret = config.Secret
	c.config.appKey = config.AppKey
}

// SetClientIP 配置
func (c *Client) SetClientIP(clientIP string) {
	if clientIP == "" {
		return
	}
	c.config.clientIP = clientIP
	if c.httpClient != nil {
		c.httpClient.SetClientIP(clientIP)
	}
}

// SetLogFun 设置日志记录函数
func (c *Client) SetLogFun(logFun gorequest.LogFunc) {
	if c.httpClient != nil {
		c.httpClient.SetLogFunc(logFun)
	}
}
