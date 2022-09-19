package dorm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne 插入一个文档
func (csc *MongoSessionCollectionOptions) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return csc.dbCollection.InsertOne(csc.session, document, opts...)
}

// InsertMany 插入多个文档
func (csc *MongoSessionCollectionOptions) InsertMany(document []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return csc.dbCollection.InsertMany(csc.session, document, opts...)
}

// DeleteOne 删除一个文档
func (csc *MongoSessionCollectionOptions) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return csc.dbCollection.DeleteOne(csc.session, filter, opts...)
}

// DeleteMany 删除多个文档
func (csc *MongoSessionCollectionOptions) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return csc.dbCollection.DeleteMany(csc.session, filter, opts...)
}

// UpdateOne 更新一个文档
func (csc *MongoSessionCollectionOptions) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return csc.dbCollection.UpdateOne(csc.session, filter, update, opts...)
}

// UpdateMany 更新多个文档
func (csc *MongoSessionCollectionOptions) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return csc.dbCollection.UpdateMany(csc.session, filter, update, opts...)
}

// FindOne 查询一个文档
func (csc *MongoSessionCollectionOptions) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return csc.dbCollection.FindOne(csc.session, filter, opts...)
}

// Find 查询多个文档
func (csc *MongoSessionCollectionOptions) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return csc.dbCollection.Find(csc.session, filter, opts...)
}
