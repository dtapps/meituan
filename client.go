package meituan

import (
	"go.dtapp.net/dorm"
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
	"gorm.io/gorm"
)

type ConfigClient struct {
	Secret       string            // 秘钥
	AppKey       string            // 渠道标记
	MongoDb      *dorm.MongoClient // 日志数据库
	PgsqlDb      *gorm.DB          // 日志数据库
	DatabaseName string            // 库名
}

// Client 美团联盟
type Client struct {
	client *gorequest.App   // 请求客户端
	log    *golog.ApiClient // 日志服务
	config *ConfigClient    // 配置
}

func NewClient(config *ConfigClient) (*Client, error) {

	var err error
	c := &Client{config: config}

	c.client = gorequest.NewHttp()

	if c.config.PgsqlDb != nil {
		c.log, err = golog.NewApiClient(
			golog.WithGormClient(c.config.PgsqlDb),
			golog.WithTableName(logTable),
		)
		if err != nil {
			return nil, err
		}
	}
	if c.config.MongoDb != nil {
		c.log, err = golog.NewApiClient(
			golog.WithMongoClient(c.config.MongoDb),
			golog.WithDatabaseName(c.config.DatabaseName),
			golog.WithCollectionName(logTable),
		)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
