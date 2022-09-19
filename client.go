package meituan

import (
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
)

// ClientConfig 实例配置
type ClientConfig struct {
	Secret           string             // 秘钥
	AppKey           string             // 渠道标记
	apiGormClientFun golog.ApiClientFun // 日志配置
	Debug            bool               // 日志开关
	ZapLog           *golog.ZapLog      // 日志服务
	CurrentIp        string             // 当前ip
}

// Client 实例
type Client struct {
	requestClient *gorequest.App // 请求服务
	zapLog        *golog.ZapLog  // 日志服务
	currentIp     string         // 当前ip
	config        struct {
		secret string // 秘钥
		appKey string // 渠道标记
	}
	log struct {
		gorm   bool             // 日志开关
		client *golog.ApiClient // 日志服务
	}
}

// NewClient 创建实例化
func NewClient(config *ClientConfig) (*Client, error) {

	c := &Client{}

	c.zapLog = config.ZapLog

	c.currentIp = config.CurrentIp

	c.config.secret = config.Secret
	c.config.appKey = config.AppKey

	c.requestClient = gorequest.NewHttp()

	apiGormClient := config.apiGormClientFun()
	if apiGormClient != nil {
		c.log.client = apiGormClient
		c.log.gorm = true
	}

	return c, nil
}
