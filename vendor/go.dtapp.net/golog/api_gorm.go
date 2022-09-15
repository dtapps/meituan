package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"gorm.io/gorm"
	"os"
	"runtime"
	"time"
	"unicode/utf8"
)

// ApiGormClientConfig 接口实例配置
type ApiGormClientConfig struct {
	GormClientFun apiGormClientFun // 日志配置
	Debug         bool             // 日志开关
	ZapLog        *ZapLog          // 日志服务
	CurrentIp     string           // 当前ip
}

// NewApiGormClient 创建接口实例化
// client 数据库服务
// tableName 表名
func NewApiGormClient(config *ApiGormClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	if config.CurrentIp == "" {
		config.CurrentIp = goip.GetOutsideIp(ctx)
	}
	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.currentIp = config.CurrentIp
	}

	client, tableName := config.GormClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.gormClient = client

	if tableName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.gormConfig.tableName = tableName

	err := c.gormAutoMigrate()
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.gormConfig.hostname = hostname
	c.gormConfig.insideIp = goip.GetInsideIp(ctx)
	c.gormConfig.goVersion = runtime.Version()

	c.log.gorm = true

	return c, nil
}

// 创建模型
func (c *ApiClient) gormAutoMigrate() (err error) {
	err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLog{})
	if err != nil {
		c.zapLog.WithLogger().Sugar().Infof("[golog.api.gormAutoMigrate]：%s", err)
	}
	return
}

// 模型结构体
type apiPostgresqlLog struct {
	LogId                 uint      `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`            //【记录】编号
	TraceId               string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`             //【系统】跟踪编号
	RequestTime           time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`           //【请求】时间
	RequestUri            string    `gorm:"comment:【请求】链接" json:"request_uri,omitempty"`                  //【请求】链接
	RequestUrl            string    `gorm:"comment:【请求】链接" json:"request_url,omitempty"`                  //【请求】链接
	RequestApi            string    `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`            //【请求】接口
	RequestMethod         string    `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`         //【请求】方式
	RequestParams         string    `gorm:"comment:【请求】参数" json:"request_params,omitempty"`               //【请求】参数
	RequestHeader         string    `gorm:"comment:【请求】头部" json:"request_header,omitempty"`               //【请求】头部
	RequestIp             string    `gorm:"index;comment:【请求】请求Ip" json:"request_ip,omitempty"`           //【请求】请求Ip
	ResponseHeader        string    `gorm:"comment:【返回】头部" json:"response_header,omitempty"`              //【返回】头部
	ResponseStatusCode    int       `gorm:"index;comment:【返回】状态码" json:"response_status_code,omitempty"`  //【返回】状态码
	ResponseBody          string    `gorm:"comment:【返回】数据" json:"response_content,omitempty"`             //【返回】数据
	ResponseContentLength int64     `gorm:"comment:【返回】大小" json:"response_content_length,omitempty"`      //【返回】大小
	ResponseTime          time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`          //【返回】时间
	SystemHostName        string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`      //【系统】主机名
	SystemInsideIp        string    `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"`     //【系统】内网ip
	SystemOs              string    `gorm:"index;comment:【系统】系统类型" json:"system_os,omitempty"`            //【系统】系统类型
	SystemArch            string    `gorm:"index;comment:【系统】系统架构" json:"system_arch,omitempty"`          //【系统】系统架构
	SystemCpuQuantity     int       `gorm:"index;comment:【系统】CPU核数" json:"system_cpu_quantity,omitempty"` //【系统】CPU核数
	GoVersion             string    `gorm:"index;comment:【程序】Go版本" json:"go_version,omitempty"`           //【程序】Go版本
	SdkVersion            string    `gorm:"index;comment:【程序】Sdk版本" json:"sdk_version,omitempty"`         //【程序】Sdk版本
}

// 记录日志
func (c *ApiClient) gormRecord(ctx context.Context, postgresqlLog apiPostgresqlLog) (err error) {

	if utf8.ValidString(postgresqlLog.ResponseBody) == false {
		postgresqlLog.ResponseBody = ""
	}

	postgresqlLog.SystemHostName = c.gormConfig.hostname
	postgresqlLog.SystemInsideIp = c.gormConfig.insideIp
	postgresqlLog.GoVersion = c.gormConfig.goVersion

	postgresqlLog.TraceId = gotrace_id.GetTraceIdContext(ctx)

	postgresqlLog.RequestIp = c.currentIp

	postgresqlLog.SystemOs = c.config.os
	postgresqlLog.SystemArch = c.config.arch
	postgresqlLog.SystemCpuQuantity = c.config.maxProCs

	err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&postgresqlLog).Error
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.gormRecord]：%s", err)
	}

	return
}

// GormQuery 查询
func (c *ApiClient) GormQuery(ctx context.Context) *gorm.DB {
	return c.gormClient.Db.Table(c.gormConfig.tableName)
}

// GormDelete 删除
func (c *ApiClient) GormDelete(ctx context.Context, hour int64) error {
	return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&apiPostgresqlLog{}).Error
}

// GormMiddleware 中间件
func (c *ApiClient) GormMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                            //【请求】时间
		RequestUri:            request.RequestUri,                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,         //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,        //【请求】接口
		RequestMethod:         request.RequestMethod,                          //【请求】方式
		RequestParams:         dorm.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:         dorm.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:        dorm.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                  //【返回】大小
		ResponseTime:          request.ResponseTime,                           //【返回】时间
		SdkVersion:            sdkVersion,                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddleware.isimg]：%s，%s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonEncodeNoError(dorm.JsonDecodeNoError(request.ResponseBody)) //【返回】数据
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddleware.len]：%s，%s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddleware.data]：%+v", data)
	}

	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.GormMiddleware]：%s", err.Error())
	}
}

// GormMiddlewareXml 中间件
func (c *ApiClient) GormMiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                            //【请求】时间
		RequestUri:            request.RequestUri,                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,         //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,        //【请求】接口
		RequestMethod:         request.RequestMethod,                          //【请求】方式
		RequestParams:         dorm.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:         dorm.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:        dorm.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                  //【返回】大小
		ResponseTime:          request.ResponseTime,                           //【返回】时间
		SdkVersion:            sdkVersion,                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareXml.isimg]：%s，%s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonEncodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareXml.len]：%s，%s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareXml.data]：%+v", data)
	}

	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.GormMiddlewareXml]：%s", err.Error())
	}
}

// GormMiddlewareCustom 中间件
func (c *ApiClient) GormMiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                            //【请求】时间
		RequestUri:            request.RequestUri,                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,         //【请求】链接
		RequestApi:            api,                                            //【请求】接口
		RequestMethod:         request.RequestMethod,                          //【请求】方式
		RequestParams:         dorm.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:         dorm.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:        dorm.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                  //【返回】大小
		ResponseTime:          request.ResponseTime,                           //【返回】时间
		SdkVersion:            sdkVersion,                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareCustom.isimg]：%s，%s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonEncodeNoError(dorm.JsonDecodeNoError(request.ResponseBody)) //【返回】数据
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareCustom.len]：%s，%s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.GormMiddlewareCustom.data]：%+v", data)
	}

	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.GormMiddlewareCustom]：%s", err.Error())
	}
}
