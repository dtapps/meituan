package meituan

import (
	"go.dtapp.net/dorm"
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
)

// Client 实例
type Client struct {
	requestClient *gorequest.App // 请求服务
	config        struct {
		secret string // 秘钥
		appKey string // 渠道标记
	}
	log struct {
		gormClient     *dorm.GormClient  // 日志数据库
		gorm           bool              // 日志开关
		logGormClient  *golog.ApiClient  // 日志服务
		mongoClient    *dorm.MongoClient // 日志数据库
		mongo          bool              // 日志开关
		logMongoClient *golog.ApiClient  // 日志服务
	}
}

// client *dorm.GormClient
type gormClientFun func() *dorm.GormClient

// client *dorm.MongoClient
// databaseName string
type mongoClientFun func() (*dorm.MongoClient, string)

// NewClient 创建实例化
// secret 秘钥
// appKey 渠道标记
func NewClient(secret, appKey string, gormClientFun gormClientFun, mongoClientFun mongoClientFun, debug bool) (*Client, error) {

	var err error
	c := &Client{}

	c.config.secret = secret
	c.config.appKey = appKey

	c.requestClient = gorequest.NewHttp()

	gormClient := gormClientFun()
	if gormClient.Db != nil {
		c.log.logGormClient, err = golog.NewApiGormClient(func() (*dorm.GormClient, string) {
			return gormClient, logTable
		}, debug)
		if err != nil {
			return nil, err
		}
		c.log.gorm = true
	}
	c.log.gormClient = gormClient

	mongoClient, databaseName := mongoClientFun()
	if mongoClient.Db != nil {
		c.log.logMongoClient, err = golog.NewApiMongoClient(func() (*dorm.MongoClient, string, string) {
			return mongoClient, databaseName, logTable
		}, debug)
		if err != nil {
			return nil, err
		}
		c.log.mongo = true
	}
	c.log.mongoClient = mongoClient

	return c, nil
}
