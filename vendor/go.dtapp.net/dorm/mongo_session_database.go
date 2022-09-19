package dorm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSessionDatabaseOptions struct {
	session    mongo.SessionContext // 会话
	dbDatabase *mongo.Database      // 数据库
}

// Database 选择数据库
func (cs *MongoSessionOptions) Database(name string, opts ...*options.DatabaseOptions) *MongoSessionDatabaseOptions {
	return &MongoSessionDatabaseOptions{
		session:    cs.Session,                    // 会话
		dbDatabase: cs.Db.Database(name, opts...), // 数据库
	}
}
