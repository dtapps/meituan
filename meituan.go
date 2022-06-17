package meituan

import (
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
	"gorm.io/gorm"
)

type ConfigClient struct {
	Secret  string   // 秘钥
	AppKey  string   // 渠道标记
	PgsqlDb *gorm.DB // pgsql数据库
}

// Client 美团联盟
type Client struct {
	client       *gorequest.App // 请求客户端
	log          *golog.Api     // 日志服务
	logTableName string         // 日志表名
	logStatus    bool           // 日志状态
	config       *ConfigClient
}

func NewClient(config *ConfigClient) *Client {

	c := &Client{}
	c.config = config

	c.client = gorequest.NewHttp()
	if c.config.PgsqlDb != nil {
		c.logStatus = true
		c.logTableName = "meituan"
		c.log = golog.NewApi(&golog.ApiConfig{
			Db:        c.config.PgsqlDb,
			TableName: c.logTableName,
		})
	}

	return c
}

func (c *Client) request(url string, params map[string]interface{}, method string) (resp gorequest.Response, err error) {

	// 创建请求
	client := c.client

	// 设置请求地址
	client.SetUri(url)

	// 设置方式
	client.SetMethod(method)

	// 设置格式
	client.SetContentTypeJson()

	// 设置参数
	client.SetParams(params)

	// 发起请求
	request, err := client.Request()
	if err != nil {
		return gorequest.Response{}, err
	}

	// 日志
	if c.logStatus == true {
		go c.postgresqlLog(request)
	}

	return request, err
}

func (c *Client) GetAppKey() string {
	return c.config.AppKey
}
