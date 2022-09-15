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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"runtime"
)

// ApiMongoClientConfig 接口实例配置
type ApiMongoClientConfig struct {
	MongoClientFun apiMongoClientFun // 日志配置
	Debug          bool              // 日志开关
	ZapLog         *ZapLog           // 日志服务
	CurrentIp      string            // 当前ip
}

// NewApiMongoClient 创建接口实例化
// client 数据库服务
// databaseName 库名
// collectionName 表名
func NewApiMongoClient(config *ApiMongoClientConfig) (*ApiClient, error) {

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

	client, databaseName, collectionName := config.MongoClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.mongoClient = client

	if databaseName == "" {
		return nil, errors.New("没有设置库名")
	}
	c.mongoConfig.databaseName = databaseName

	if collectionName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.mongoConfig.collectionName = collectionName

	hostname, _ := os.Hostname()

	c.mongoConfig.hostname = hostname
	c.mongoConfig.insideIp = goip.GetInsideIp(ctx)
	c.mongoConfig.goVersion = runtime.Version()

	c.log.mongo = true

	// 创建时间序列集合
	c.mongoCreateCollection(ctx)

	// 创建索引
	c.mongoCreateIndexes(ctx)

	return c, nil
}

// 创建时间序列集合
func (c *ApiClient) mongoCreateCollection(ctx context.Context) {
	var commandResult bson.M
	commandErr := c.mongoClient.Db.Database(c.mongoConfig.databaseName).RunCommand(ctx, bson.D{{
		"listCollections", 1,
	}}).Decode(&commandResult)
	if commandErr != nil {
		c.zapLog.WithLogger().Sugar().Errorf("检查时间序列集合：%s", commandErr)
	} else {
		err := c.mongoClient.Db.Database(c.mongoConfig.databaseName).CreateCollection(ctx, c.mongoConfig.collectionName, options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
		if err != nil {
			c.zapLog.WithLogger().Sugar().Errorf("创建时间序列集合：%s", err)
		}
	}
}

// 创建索引
func (c *ApiClient) mongoCreateIndexes(ctx context.Context) {
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"trace_id", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"log_time", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"request_time", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"request_method", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"response_status_code", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"response_time", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"system_host_name", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"system_inside_ip", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"system_os", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"system_arch", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"system_cpu_quantity", 1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"go_version", -1},
	}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"sdk_version", -1},
	}}))
}

// 模型结构体
type apiMongolLog struct {
	LogId                 primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                                      //【记录】编号
	LogTime               primitive.DateTime `json:"log_time,omitempty" bson:"log_time,omitempty"`                               //【记录】时间
	TraceId               string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                               //【记录】跟踪编号
	RequestTime           dorm.BsonTime      `json:"request_time,omitempty" bson:"request_time,omitempty"`                       //【请求】时间
	RequestUri            string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`                         //【请求】链接
	RequestUrl            string             `json:"request_url,omitempty" bson:"request_url,omitempty"`                         //【请求】链接
	RequestApi            string             `json:"request_api,omitempty" bson:"request_api,omitempty"`                         //【请求】接口
	RequestMethod         string             `json:"request_method,omitempty" bson:"request_method,omitempty"`                   //【请求】方式
	RequestParams         interface{}        `json:"request_params,omitempty" bson:"request_params,omitempty"`                   //【请求】参数
	RequestHeader         interface{}        `json:"request_header,omitempty" bson:"request_header,omitempty"`                   //【请求】头部
	RequestIp             string             `json:"request_ip,omitempty" bson:"request_ip,omitempty"`                           //【请求】请求Ip
	ResponseHeader        interface{}        `json:"response_header,omitempty" bson:"response_header,omitempty"`                 //【返回】头部
	ResponseStatusCode    int                `json:"response_status_code,omitempty" bson:"response_status_code,omitempty"`       //【返回】状态码
	ResponseBody          interface{}        `json:"response_body,omitempty" bson:"response_body,omitempty"`                     //【返回】内容
	ResponseContentLength int64              `json:"response_content_length,omitempty" bson:"response_content_length,omitempty"` //【返回】大小
	ResponseTime          dorm.BsonTime      `json:"response_time,omitempty" bson:"response_time,omitempty"`                     //【返回】时间
	SystemHostName        string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`               //【系统】主机名
	SystemInsideIp        string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`               //【系统】内网ip
	SystemOs              string             `json:"system_os,omitempty" bson:"system_os,omitempty"`                             //【系统】系统类型
	SystemArch            string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                         //【系统】系统架构
	SystemCpuQuantity     int                `json:"system_cpu_quantity,omitempty" bson:"system_cpu_quantity,omitempty"`         //【系统】CPU核数
	GoVersion             string             `json:"go_version,omitempty" bson:"go_version,omitempty"`                           //【程序】Go版本
	SdkVersion            string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                         //【程序】Sdk版本
}

// 记录日志
func (c *ApiClient) mongoRecord(ctx context.Context, mongoLog apiMongolLog) (err error) {

	mongoLog.SystemHostName = c.mongoConfig.hostname     //【系统】主机名
	mongoLog.SystemInsideIp = c.mongoConfig.insideIp     //【系统】内网ip
	mongoLog.GoVersion = c.mongoConfig.goVersion         //【程序】Go版本
	mongoLog.TraceId = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	mongoLog.RequestIp = c.currentIp                     //【请求】请求Ip
	mongoLog.SystemOs = c.config.os                      //【系统】系统类型
	mongoLog.SystemArch = c.config.arch                  //【系统】系统架构
	mongoLog.SystemCpuQuantity = c.config.maxProCs       //【系统】CPU核数
	mongoLog.LogId = primitive.NewObjectID()             //【记录】编号

	_, err = c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(mongoLog)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.mongoRecord]：%s", err)
	}

	return err
}

// MongoQuery 查询
func (c *ApiClient) MongoQuery(ctx context.Context) *mongo.Collection {
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName)
}

// MongoDelete 删除
func (c *ApiClient) MongoDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).DeleteMany(ctx, filter)
}

// MongoMiddleware 中间件
func (c *ApiClient) MongoMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.BsonTime(request.RequestTime),                 //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.BsonTime(request.ResponseTime),                //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddleware.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddleware.len]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddleware.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.MongoMiddleware]：%s", err.Error())
	}
}

// MongoMiddlewareXml 中间件
func (c *ApiClient) MongoMiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.BsonTime(request.RequestTime),                 //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.BsonTime(request.ResponseTime),                //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddlewareXml.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.XmlDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddlewareXml]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddlewareXml.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.MongoMiddlewareXml]：%s", err.Error())
	}
}

// MongoMiddlewareCustom 中间件
func (c *ApiClient) MongoMiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.BsonTime(request.RequestTime),                 //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            api,                                                //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.BsonTime(request.ResponseTime),                //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddlewareCustom.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.MongoMiddlewareCustom]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoRecordJson.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.MongoMiddlewareCustom]：%s", err.Error())
	}
}
