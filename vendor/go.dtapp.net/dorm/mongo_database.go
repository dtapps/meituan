package dorm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseOptions struct {
	dbDatabase *mongo.Database // 数据库
}

// Database 选择数据库
func (c *MongoClient) Database(name string, opts ...*options.DatabaseOptions) *MongoDatabaseOptions {
	return &MongoDatabaseOptions{
		dbDatabase: c.Db.Database(name, opts...),
	}
}

// CreateCollection 创建集合
func (cd *MongoDatabaseOptions) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	return cd.dbDatabase.CreateCollection(ctx, name, opts...)
}

// CreateTimeSeriesCollection 创建时间序列集合
func (cd *MongoDatabaseOptions) CreateTimeSeriesCollection(ctx context.Context, name string, timeField string) error {
	return cd.dbDatabase.CreateCollection(ctx, name, options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField(timeField)))
}
