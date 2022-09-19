package dorm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClientFun *MongoClient 驱动
// string 库名
type MongoClientFun func() (*MongoClient, string)

// MongoClientCollectionFun *MongoClient 驱动
// string 库名
// string 集合
type MongoClientCollectionFun func() (*MongoClient, string, string)

type ConfigMongoClient struct {
	Dns          string // 地址
	Opts         *options.ClientOptions
	DatabaseName string // 库名
}

type MongoClient struct {
	Db     *mongo.Client      // 驱动
	config *ConfigMongoClient // 配置
}

func NewMongoClient(config *ConfigMongoClient) (*MongoClient, error) {

	var ctx = context.Background()
	var err error
	c := &MongoClient{config: config}

	// 连接到MongoDB
	if c.config.Dns != "" {
		c.Db, err = mongo.Connect(ctx, options.Client().ApplyURI(c.config.Dns))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("连接失败：%v", err))
		}
	} else {
		c.Db, err = mongo.Connect(ctx, c.config.Opts)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("连接失败：%v", err))
		}
	}

	// 检查连接
	err = c.Db.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("检查连接失败：%v", err))
	}

	return c, nil
}

// Close 关闭
func (c *MongoClient) Close(ctx context.Context) error {
	return c.Db.Disconnect(ctx)
}
