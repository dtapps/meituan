package dorm

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSessionOptions struct {
	Session      mongo.SessionContext // 会话
	Db           *mongo.Client        // 驱动
	startSession mongo.Session        // 开始会话
}

// Begin 开始事务，会同时创建开始会话需要在退出时关闭会话
func (c *MongoClient) Begin() (ms *MongoSessionOptions, err error) {

	var ctx = context.Background()

	ms.Db = c.Db

	// 开始会话
	ms.startSession, err = ms.Db.StartSession()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("开始会话失败：%v", err))
	}

	// 会话上下文
	ms.Session = mongo.NewSessionContext(ctx, ms.startSession)

	// 会话开启事务
	err = ms.startSession.StartTransaction()
	return ms, err
}

// Close 关闭会话
func (cs *MongoSessionOptions) Close(ctx context.Context) {
	cs.startSession.EndSession(ctx)
}

// Rollback 回滚事务
func (cs *MongoSessionOptions) Rollback(ctx context.Context) error {
	return cs.startSession.AbortTransaction(ctx)
}

// Commit 提交事务
func (cs *MongoSessionOptions) Commit(ctx context.Context) error {
	return cs.startSession.CommitTransaction(ctx)
}
