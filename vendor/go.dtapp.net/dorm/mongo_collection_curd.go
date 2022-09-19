package dorm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne 插入一个文档
func (cc *MongoCollectionOptions) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return cc.dbCollection.InsertOne(ctx, document, opts...)
}

// InsertMany 插入多个文档
func (cc *MongoCollectionOptions) InsertMany(ctx context.Context, document []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return cc.dbCollection.InsertMany(ctx, document, opts...)
}

// DeleteOne 删除一个文档
func (cc *MongoCollectionOptions) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return cc.dbCollection.DeleteOne(ctx, filter, opts...)
}

// DeleteMany 删除多个文档
func (cc *MongoCollectionOptions) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return cc.dbCollection.DeleteMany(ctx, filter, opts...)
}

// UpdateOne 更新一个文档
func (cc *MongoCollectionOptions) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return cc.dbCollection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany 更新多个文档
func (cc *MongoCollectionOptions) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return cc.dbCollection.UpdateMany(ctx, filter, update, opts...)
}

// FindOne 查询一个文档
func (cc *MongoCollectionOptions) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return cc.dbCollection.FindOne(ctx, filter, opts...)
}

// Find 查询多个文档
func (cc *MongoCollectionOptions) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return cc.dbCollection.Find(ctx, filter, opts...)
}
