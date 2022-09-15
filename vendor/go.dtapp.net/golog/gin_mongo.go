package golog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
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
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"time"
)

// GinMongoClientConfig 框架实例配置
type GinMongoClientConfig struct {
	IpService      *goip.Client      // ip服务
	MongoClientFun ginMongoClientFun // 日志配置
	Debug          bool              // 日志开关
	ZapLog         *ZapLog           // 日志服务
}

// NewGinMongoClient 创建框架实例化
// client 数据库服务
// databaseName 库名
// collectionName 表名
// ipService ip服务
func NewGinMongoClient(config *GinMongoClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

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

	c.ipService = config.IpService

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
func (c *GinClient) mongoCreateCollection(ctx context.Context) {
	var commandResult bson.M
	commandErr := c.mongoClient.Db.Database(c.mongoConfig.databaseName).RunCommand(ctx, bson.D{{
		"listCollections", 1,
	}}).Decode(&commandResult)
	if commandErr != nil {
		c.zapLog.WithLogger().Sugar().Error("检查时间序列集合：", commandErr)
	} else {
		err := c.mongoClient.Db.Database(c.mongoConfig.databaseName).CreateCollection(ctx, c.mongoConfig.collectionName, options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
		if err != nil {
			c.zapLog.WithLogger().Sugar().Error("创建时间序列集合：", err)
		}
	}
}

// 创建索引
func (c *GinClient) mongoCreateIndexes(ctx context.Context) {
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"trace_id", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"log_time", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_time", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_method", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_proto", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip_country", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip_region", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip_province", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip_city", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"request_ip_isp", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"response_time", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"response_code", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"system_host_name", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"system_inside_ip", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"system_os", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"system_arch", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"system_cpu_quantity", 1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"go_version", -1},
		}}))
	c.zapLog.WithLogger().Sugar().Infof(c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"sdk_version", -1},
		}}))
}

// 模型结构体
type ginMongoLog struct {
	LogId             primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                              //【记录】编号
	LogTime           primitive.DateTime `json:"log_time,omitempty" bson:"log_time,omitempty"`                       //【记录】时间
	TraceId           string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                       //【记录】跟踪编号
	RequestTime       dorm.BsonTime      `json:"request_time,omitempty" bson:"request_time,omitempty"`               //【请求】时间
	RequestUri        string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`                 //【请求】请求链接 域名+路径+参数
	RequestUrl        string             `json:"request_url,omitempty" bson:"request_url,omitempty"`                 //【请求】请求链接 域名+路径
	RequestApi        string             `json:"request_api,omitempty" bson:"request_api,omitempty"`                 //【请求】请求接口 路径
	RequestMethod     string             `json:"request_method,omitempty" bson:"request_method,omitempty"`           //【请求】请求方式
	RequestProto      string             `json:"request_proto,omitempty" bson:"request_proto,omitempty"`             //【请求】请求协议
	RequestUa         string             `json:"request_ua,omitempty" bson:"request_ua,omitempty"`                   //【请求】请求UA
	RequestReferer    string             `json:"request_referer,omitempty" bson:"request_referer,omitempty"`         //【请求】请求referer
	RequestBody       interface{}        `json:"request_body,omitempty" bson:"request_body,omitempty"`               //【请求】请求主体
	RequestUrlQuery   interface{}        `json:"request_url_query,omitempty" bson:"request_url_query,omitempty"`     //【请求】请求URL参数
	RequestIp         string             `json:"request_ip,omitempty" bson:"request_ip,omitempty"`                   //【请求】请求客户端Ip
	RequestIpCountry  string             `json:"request_ip_country,omitempty" bson:"request_ip_country,omitempty"`   //【请求】请求客户端城市
	RequestIpRegion   string             `json:"request_ip_region,omitempty" bson:"request_ip_region,omitempty"`     //【请求】请求客户端区域
	RequestIpProvince string             `json:"request_ip_province,omitempty" bson:"request_ip_province,omitempty"` //【请求】请求客户端省份
	RequestIpCity     string             `json:"request_ip_city,omitempty" bson:"request_ip_city,omitempty"`         //【请求】请求客户端城市
	RequestIpIsp      string             `json:"request_ip_isp,omitempty" bson:"request_ip_isp,omitempty"`           //【请求】请求客户端运营商
	RequestHeader     interface{}        `json:"request_header,omitempty" bson:"request_header,omitempty"`           //【请求】请求头
	ResponseTime      dorm.BsonTime      `json:"response_time,omitempty" bson:"response_time,omitempty"`             //【返回】时间
	ResponseCode      int                `json:"response_code,omitempty" bson:"response_code,omitempty"`             //【返回】状态码
	ResponseMsg       string             `json:"response_msg,omitempty" bson:"response_msg,omitempty"`               //【返回】描述
	ResponseData      interface{}        `json:"response_data,omitempty" bson:"response_data,omitempty"`             //【返回】数据
	CostTime          int64              `json:"cost_time,omitempty" bson:"cost_time,omitempty"`                     //【系统】花费时间
	SystemHostName    string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`       //【系统】主机名
	SystemInsideIp    string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`       //【系统】内网ip
	SystemOs          string             `json:"system_os,omitempty" bson:"system_os,omitempty"`                     //【系统】系统类型
	SystemArch        string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                 //【系统】系统架构
	SystemCpuQuantity int                `json:"system_cpu_quantity,omitempty" bson:"system_cpu_quantity,omitempty"` //【系统】CPU核数
	GoVersion         string             `json:"go_version,omitempty" bson:"go_version,omitempty"`                   //【程序】Go版本
	SdkVersion        string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                 //【程序】Sdk版本
}

// 记录日志
func (c *GinClient) mongoRecord(mongoLog ginMongoLog) (err error) {

	mongoLog.SystemHostName = c.mongoConfig.hostname //【系统】主机名
	mongoLog.SystemInsideIp = c.mongoConfig.insideIp //【系统】内网ip
	mongoLog.GoVersion = c.mongoConfig.goVersion     //【程序】Go版本
	mongoLog.SdkVersion = Version                    //【程序】Sdk版本
	mongoLog.SystemOs = c.config.os                  //【系统】系统类型
	mongoLog.SystemArch = c.config.arch              //【系统】系统架构
	mongoLog.SystemCpuQuantity = c.config.maxProCs   //【系统】CPU核数
	mongoLog.LogId = primitive.NewObjectID()         //【记录】编号

	_, err = c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(mongoLog)
	if err != nil {
		c.zapLog.WithTraceIdStr(mongoLog.TraceId).Sugar().Errorf("[golog.gin.mongoRecord]：%s", err)
	}

	return err
}

func (c *GinClient) mongoRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.mongoRecordJson]收到保存数据要求：%s,%s", c.mongoConfig.databaseName, c.mongoConfig.collectionName)
	}

	data := ginMongoLog{
		TraceId:           traceId,                                                      //【记录】跟踪编号
		LogTime:           primitive.NewDateTimeFromTime(requestTime),                   //【记录】时间
		RequestTime:       dorm.BsonTime(requestTime),                                   //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                        //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:      dorm.BsonTime(gotime.Current().Time),                         //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      c.jsonUnmarshal(responseBody),                                //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.JsonDecodeNoError(requestBody) //【请求】请求主体
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordJson.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordJson.data]：%+v", data)
	}

	err := c.mongoRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.mongoRecordJson]：%s", err)
	}
}

func (c *GinClient) mongoRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.mongoRecordXml]收到保存数据要求：%s,%s", c.mongoConfig.databaseName, c.mongoConfig.collectionName)
	}

	data := ginMongoLog{
		TraceId:           traceId,                                                      //【记录】跟踪编号
		LogTime:           primitive.NewDateTimeFromTime(requestTime),                   //【记录】时间
		RequestTime:       dorm.BsonTime(requestTime),                                   //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                        //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:      dorm.BsonTime(gotime.Current().Time),                         //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      c.jsonUnmarshal(responseBody),                                //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.XmlDecodeNoError(requestBody) //【请求】请求主体
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordXml.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordXml.data]：%+v", data)
	}

	err := c.mongoRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.mongoRecordXml]：%s", err)
	}
}

// MongoQuery 查询
func (c *GinClient) MongoQuery(ctx context.Context) *mongo.Collection {
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName)
}

// MongoDelete 删除
func (c *GinClient) MongoDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).DeleteMany(ctx, filter)
}

// MongoMiddleware 中间件
func (c *GinClient) MongoMiddleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取
		data, _ := ioutil.ReadAll(ginCtx.Request.Body)

		if c.logDebug {
			c.zapLog.WithLogger().Sugar().Infof("[golog.gin.MongoMiddleware] %s", data)
		}

		// 复用
		ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := blw.body.String()

		//结束时间
		endTime := gotime.Current().TimestampWithMillisecond()

		go func() {

			var dataJson = true

			// 解析请求内容
			var jsonBody map[string]interface{}

			// 判断是否有内容
			if len(data) > 0 {
				err := json.Unmarshal(data, &jsonBody)
				if err != nil {
					dataJson = false
				}
			}

			clientIp := gorequest.ClientIp(ginCtx.Request)

			requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp := "", "", "", "", ""
			if c.ipService != nil {
				if net.ParseIP(clientIp).To4() != nil {
					// IPv4
					_, info := c.ipService.Ipv4(clientIp)
					requestClientIpCountry = info.Country
					requestClientIpRegion = info.Region
					requestClientIpProvince = info.Province
					requestClientIpCity = info.City
					requestClientIpIsp = info.ISP
				} else if net.ParseIP(clientIp).To16() != nil {
					// IPv6
					info := c.ipService.Ipv6(clientIp)
					requestClientIpCountry = info.Country
					requestClientIpProvince = info.Province
					requestClientIpCity = info.City
				}
			}

			// 记录
			if c.mongoClient != nil && c.mongoClient.Db != nil {

				var traceId = gotrace_id.GetGinTraceId(ginCtx)

				if dataJson {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.MongoMiddleware]准备使用{mongoRecordJson}保存数据：%s", data)
					}
					c.mongoRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				} else {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.MongoMiddleware]准备使用{mongoRecordXml}保存数据：%s", data)
					}
					c.mongoRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				}
			}
		}()
	}
}
