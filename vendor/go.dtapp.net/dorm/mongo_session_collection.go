package dorm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSessionCollectionOptions struct {
	session      mongo.SessionContext // 会话
	dbCollection *mongo.Collection    // 集合
}

// Collection 选择集合
func (csd *MongoSessionDatabaseOptions) Collection(name string, opts ...*options.CollectionOptions) *MongoSessionCollectionOptions {
	return &MongoSessionCollectionOptions{
		session:      csd.session,                              // 会话
		dbCollection: csd.dbDatabase.Collection(name, opts...), // 集合
	}
}
