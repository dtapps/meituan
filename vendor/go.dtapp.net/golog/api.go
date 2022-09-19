package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"os"
	"runtime"
)

// ApiClientFun *ApiClient 驱动
type ApiClientFun func() *ApiClient

// ApiClientJsonFun *ApiClient 驱动
// jsonStatus bool json状态
type ApiClientJsonFun func() (*ApiClient, bool)

// ApiClient 接口
type ApiClient struct {
	gormClient *dorm.GormClient // 数据库驱动
	zapLog     *ZapLog          // 日志服务
	logDebug   bool             // 日志开关
	currentIp  string           // 当前ip
	gormConfig struct {
		tableName string // 表名
		insideIp  string // 内网ip
		hostname  string // 主机名
		goVersion string // go版本
	}
	config struct {
		os         string // 系统类型
		arch       string // 系统架构
		jsonStatus bool   // json状态
	}
}

// ApiClientConfig 接口实例配置
type ApiClientConfig struct {
	GormClientFun dorm.GormClientTableFun // 日志配置
	Debug         bool                    // 日志开关
	ZapLog        *ZapLog                 // 日志服务
	CurrentIp     string                  // 当前ip
	JsonStatus    bool                    // json状态
}

// NewApiClient 创建接口实例化
func NewApiClient(config *ApiClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	c.config.jsonStatus = config.JsonStatus

	if config.CurrentIp == "" {
		config.CurrentIp = goip.GetOutsideIp(ctx)
	}
	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.currentIp = config.CurrentIp
	}

	c.config.os = runtime.GOOS
	c.config.arch = runtime.GOARCH

	gormClient, gormTableName := config.GormClientFun()

	if gormClient == nil || gormClient.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	hostname, _ := os.Hostname()

	if gormClient != nil || gormClient.Db != nil {

		c.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		}
		c.gormConfig.tableName = gormTableName

		err := c.gormAutoMigrate()
		if err != nil {
			return nil, errors.New("创建表失败：" + err.Error())
		}

		c.gormConfig.hostname = hostname
		c.gormConfig.insideIp = goip.GetInsideIp(ctx)
		c.gormConfig.goVersion = runtime.Version()

	}

	return c, nil
}

// Middleware 中间件
func (c *ApiClient) Middleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	c.GormMiddleware(ctx, request, sdkVersion)
}

// MiddlewareXml 中间件
func (c *ApiClient) MiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	c.GormMiddlewareXml(ctx, request, sdkVersion)
}

// MiddlewareCustom 中间件
func (c *ApiClient) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	c.GormMiddlewareCustom(ctx, api, request, sdkVersion)
}
