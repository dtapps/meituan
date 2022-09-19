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
	"gorm.io/datatypes"
	"os"
	"runtime"
	"unicode/utf8"
)

// ApiGormClientConfig 接口实例配置
type ApiGormClientConfig struct {
	GormClientFun dorm.GormClientTableFun // 日志配置
	Debug         bool                    // 日志开关
	ZapLog        *ZapLog                 // 日志服务
	CurrentIp     string                  // 当前ip
	JsonStatus    bool                    // json状态
}

// NewApiGormClient 创建接口实例化
func NewApiGormClient(config *ApiGormClientConfig) (*ApiClient, error) {

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

	return c, nil
}

// 创建模型
func (c *ApiClient) gormAutoMigrate() (err error) {
	if c.config.jsonStatus {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLogJson{})
		if err != nil {
			c.zapLog.WithLogger().Sugar().Errorf("创建模型：%s", err)
		}
	} else {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLogString{})
		if err != nil {
			c.zapLog.WithLogger().Sugar().Errorf("创建模型：%s", err)
		}
	}
	return nil
}

// 记录日志
func (c *ApiClient) gormRecord(ctx context.Context, data apiPostgresqlLogString) (err error) {

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	data.SystemHostName = c.gormConfig.hostname
	data.SystemInsideIp = c.gormConfig.insideIp
	data.GoVersion = c.gormConfig.goVersion

	data.TraceId = gotrace_id.GetTraceIdContext(ctx)

	data.RequestIp = c.currentIp

	data.SystemOs = c.config.os
	data.SystemArch = c.config.arch

	if c.config.jsonStatus {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&apiPostgresqlLogJson{
			TraceId:               data.TraceId,
			RequestTime:           data.RequestTime,
			RequestUri:            data.RequestUri,
			RequestUrl:            data.RequestUrl,
			RequestApi:            data.RequestApi,
			RequestMethod:         data.RequestMethod,
			RequestParams:         datatypes.JSON(data.RequestParams),
			RequestHeader:         datatypes.JSON(data.RequestHeader),
			RequestIp:             data.RequestIp,
			ResponseHeader:        datatypes.JSON(data.ResponseHeader),
			ResponseStatusCode:    data.ResponseStatusCode,
			ResponseBody:          datatypes.JSON(data.ResponseBody),
			ResponseContentLength: data.ResponseContentLength,
			ResponseTime:          data.ResponseTime,
			SystemHostName:        data.SystemHostName,
			SystemInsideIp:        data.SystemInsideIp,
			SystemOs:              data.SystemOs,
			SystemArch:            data.SystemArch,
			GoVersion:             data.GoVersion,
			SdkVersion:            data.SdkVersion,
		}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("记录日志失败：%s", err)
		}
	} else {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&apiPostgresqlLogString{
			TraceId:               data.TraceId,
			RequestTime:           data.RequestTime,
			RequestUri:            data.RequestUri,
			RequestUrl:            data.RequestUrl,
			RequestApi:            data.RequestApi,
			RequestMethod:         data.RequestMethod,
			RequestParams:         data.RequestParams,
			RequestHeader:         data.RequestHeader,
			RequestIp:             data.RequestIp,
			ResponseHeader:        data.ResponseHeader,
			ResponseStatusCode:    data.ResponseStatusCode,
			ResponseBody:          data.ResponseBody,
			ResponseContentLength: data.ResponseContentLength,
			ResponseTime:          data.ResponseTime,
			SystemHostName:        data.SystemHostName,
			SystemInsideIp:        data.SystemInsideIp,
			SystemOs:              data.SystemOs,
			SystemArch:            data.SystemArch,
			GoVersion:             data.GoVersion,
			SdkVersion:            data.SdkVersion,
		}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("记录日志失败：%s", err)
		}
	}

	return
}

// GormDelete 删除
func (c *ApiClient) GormDelete(ctx context.Context, hour int64) error {
	if c.config.jsonStatus {
		return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&apiPostgresqlLogJson{}).Error
	} else {
		return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&apiPostgresqlLogString{}).Error
	}
}

// GormMiddleware 中间件
func (c *ApiClient) GormMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLogString{
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
	data := apiPostgresqlLogString{
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
	data := apiPostgresqlLogString{
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
