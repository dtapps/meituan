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
	"gorm.io/datatypes"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"time"
)

// GinGormClientConfig 框架实例配置
type GinGormClientConfig struct {
	IpService     *goip.Client            // ip服务
	GormClientFun dorm.GormClientTableFun // 日志配置
	Debug         bool                    // 日志开关
	ZapLog        *ZapLog                 // 日志服务
	JsonStatus    bool                    // json状态
}

// NewGinGormClient 创建框架实例化
// client 数据库服务
// tableName 表名
// ipService ip服务
func NewGinGormClient(config *GinGormClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	c.config.jsonStatus = config.JsonStatus

	client, tableName := config.GormClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.gormClient = client

	if tableName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.gormConfig.tableName = tableName

	c.ipService = config.IpService

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
func (c *GinClient) gormAutoMigrate() (err error) {
	if c.config.jsonStatus {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLogJson{})
		if err != nil {
			c.zapLog.WithLogger().Sugar().Errorf("创建模型：%s", err)
		}
	} else {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLogString{})
		if err != nil {
			c.zapLog.WithLogger().Sugar().Errorf("创建模型：%s", err)
		}
	}
	return err
}

// gormRecord 记录日志
func (c *GinClient) gormRecord(data ginPostgresqlLogString) (err error) {

	data.SystemHostName = c.gormConfig.hostname
	data.SystemInsideIp = c.gormConfig.insideIp
	data.GoVersion = c.gormConfig.goVersion

	data.SdkVersion = Version

	data.SystemOs = c.config.os
	data.SystemArch = c.config.arch

	if c.config.jsonStatus {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&ginPostgresqlLogJson{
			TraceId:            data.TraceId,
			RequestTime:        data.RequestTime,
			RequestUri:         data.RequestUri,
			RequestUrl:         data.RequestUrl,
			RequestApi:         data.RequestApi,
			RequestMethod:      data.RequestMethod,
			RequestProto:       data.RequestProto,
			RequestUa:          data.RequestUa,
			RequestReferer:     data.RequestReferer,
			RequestBody:        datatypes.JSON(data.RequestBody),
			RequestUrlQuery:    datatypes.JSON(data.RequestUrlQuery),
			RequestIp:          data.RequestIp,
			RequestIpCountry:   data.RequestIpCountry,
			RequestIpProvince:  data.RequestIpProvince,
			RequestIpCity:      data.RequestIpCity,
			RequestIpIsp:       data.RequestIpIsp,
			RequestIpLongitude: data.RequestIpLongitude,
			RequestIpLatitude:  data.RequestIpLatitude,
			RequestHeader:      datatypes.JSON(data.RequestHeader),
			ResponseTime:       data.ResponseTime,
			ResponseCode:       data.ResponseCode,
			ResponseMsg:        data.ResponseMsg,
			ResponseData:       datatypes.JSON(data.ResponseData),
			CostTime:           data.CostTime,
			SystemHostName:     data.SystemHostName,
			SystemInsideIp:     data.SystemInsideIp,
			SystemOs:           data.SystemOs,
			SystemArch:         data.SystemArch,
			GoVersion:          data.GoVersion,
			SdkVersion:         data.SdkVersion,
		}).Error
		if err != nil {
			c.zapLog.WithTraceIdStr(data.TraceId).Sugar().Errorf("记录日志失败：%s", err)
		}
	} else {
		err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&ginPostgresqlLogString{
			TraceId:            data.TraceId,
			RequestTime:        data.RequestTime,
			RequestUri:         data.RequestUri,
			RequestUrl:         data.RequestUrl,
			RequestApi:         data.RequestApi,
			RequestMethod:      data.RequestMethod,
			RequestProto:       data.RequestProto,
			RequestUa:          data.RequestUa,
			RequestReferer:     data.RequestReferer,
			RequestBody:        data.RequestBody,
			RequestUrlQuery:    data.RequestUrlQuery,
			RequestIp:          data.RequestIp,
			RequestIpCountry:   data.RequestIpCountry,
			RequestIpProvince:  data.RequestIpProvince,
			RequestIpCity:      data.RequestIpCity,
			RequestIpIsp:       data.RequestIpIsp,
			RequestIpLongitude: data.RequestIpLongitude,
			RequestIpLatitude:  data.RequestIpLatitude,
			RequestHeader:      data.RequestHeader,
			ResponseTime:       data.ResponseTime,
			ResponseCode:       data.ResponseCode,
			ResponseMsg:        data.ResponseMsg,
			ResponseData:       data.ResponseData,
			CostTime:           data.CostTime,
			SystemHostName:     data.SystemHostName,
			SystemInsideIp:     data.SystemInsideIp,
			SystemOs:           data.SystemOs,
			SystemArch:         data.SystemArch,
			GoVersion:          data.GoVersion,
			SdkVersion:         data.SdkVersion,
		}).Error
		if err != nil {
			c.zapLog.WithTraceIdStr(data.TraceId).Sugar().Errorf("记录日志失败：%s", err)
		}
	}

	return
}

func (c *GinClient) gormRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string, requestClientIpLocationLatitude, requestClientIpLocationLongitude float64) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.gormRecordJson]收到保存数据要求：%s", c.gormConfig.tableName)
	}

	data := ginPostgresqlLogString{
		TraceId:            traceId,                                                      //【系统】跟踪编号
		RequestTime:        requestTime,                                                  //【请求】时间
		RequestUrl:         ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:         gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:      ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:       ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:          ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:     ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:    dorm.JsonEncodeNoError(ginCtx.Request.URL.Query()),           //【请求】请求URL参数
		RequestIp:          clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:   requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpProvince:  requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:      requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:       requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestIpLatitude:  requestClientIpLocationLatitude,                              // 【请求】请求客户端纬度
		RequestIpLongitude: requestClientIpLocationLongitude,                             // 【请求】请求客户端经度
		RequestHeader:      dorm.JsonEncodeNoError(ginCtx.Request.Header),                //【请求】请求头
		ResponseTime:       gotime.Current().Time,                                        //【返回】时间
		ResponseCode:       responseCode,                                                 //【返回】状态码
		ResponseData:       responseBody,                                                 //【返回】数据
		CostTime:           endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.JsonEncodeNoError(requestBody) //【请求】请求主体
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.gormRecordJson.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.gormRecordJson.data]：%+v", data)
	}

	err := c.gormRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordJson]：%s", err)
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordJson.string]：%s", requestBody)
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordJson.JsonEncodeNoError.string]：%s", dorm.JsonEncodeNoError(requestBody))
	}
}

func (c *GinClient) gormRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string, requestClientIpLocationLatitude, requestClientIpLocationLongitude float64) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.gormRecordXml]收到保存数据要求：%s", c.gormConfig.tableName)
	}

	data := ginPostgresqlLogString{
		TraceId:            traceId,                                                      //【系统】跟踪编号
		RequestTime:        requestTime,                                                  //【请求】时间
		RequestUrl:         ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:         gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:      ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:       ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:          ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:     ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:    dorm.JsonEncodeNoError(ginCtx.Request.URL.Query()),           //【请求】请求URL参数
		RequestIp:          clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:   requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpProvince:  requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:      requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:       requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestIpLatitude:  requestClientIpLocationLatitude,                              // 【请求】请求客户端纬度
		RequestIpLongitude: requestClientIpLocationLongitude,                             // 【请求】请求客户端经度
		RequestHeader:      dorm.JsonEncodeNoError(ginCtx.Request.Header),                //【请求】请求头
		ResponseTime:       gotime.Current().Time,                                        //【返回】时间
		ResponseCode:       responseCode,                                                 //【返回】状态码
		ResponseData:       responseBody,                                                 //【返回】数据
		CostTime:           endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.XmlEncodeNoError(dorm.XmlDecodeNoError(requestBody)) //【请求】请求内容
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.gormRecordXml.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.gormRecordXml.data]：%+v", data)
	}

	err := c.gormRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordXml]：%s", err)
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordXml.string]：%s", requestBody)
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordXml.XmlDecodeNoError.string]：%s", dorm.XmlDecodeNoError(requestBody))
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.gormRecordXml.XmlEncodeNoError.string]：%s", dorm.XmlEncodeNoError(dorm.XmlDecodeNoError(requestBody)))
	}
}

// GormDelete 删除
func (c *GinClient) GormDelete(ctx context.Context, hour int64) error {
	if c.config.jsonStatus {
		return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&ginPostgresqlLogJson{}).Error
	} else {
		return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&ginPostgresqlLogString{}).Error
	}
}

// GormMiddleware 中间件
func (c *GinClient) GormMiddleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取
		data, _ := ioutil.ReadAll(ginCtx.Request.Body)

		if c.logDebug {
			c.zapLog.WithLogger().Sugar().Infof("[golog.gin.GormMiddleware]：%s", data)
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

			var requestClientIpCountry string
			var requestClientIpProvince string
			var requestClientIpCity string
			var requestClientIpIsp string
			var requestClientIpLocationLatitude float64
			var requestClientIpLocationLongitude float64
			if c.ipService != nil {
				if net.ParseIP(clientIp).To4() != nil {
					// IPv4
					info := c.ipService.Analyse(clientIp)
					requestClientIpCountry = info.Ip2regionV2info.Country
					requestClientIpProvince = info.Ip2regionV2info.Province
					requestClientIpCity = info.Ip2regionV2info.City
					requestClientIpIsp = info.Ip2regionV2info.Operator
					requestClientIpLocationLatitude = info.GeoipInfo.Location.Latitude
					requestClientIpLocationLongitude = info.GeoipInfo.Location.Longitude
				} else if net.ParseIP(clientIp).To16() != nil {
					// IPv6
					info := c.ipService.Analyse(clientIp)
					requestClientIpCountry = info.Ipv6wryInfo.Country
					requestClientIpProvince = info.Ipv6wryInfo.Province
					requestClientIpCity = info.Ipv6wryInfo.City
					requestClientIpLocationLatitude = info.GeoipInfo.Location.Latitude
					requestClientIpLocationLongitude = info.GeoipInfo.Location.Longitude
				}
			}

			// 记录
			if c.gormClient != nil && c.gormClient.Db != nil {

				var traceId = gotrace_id.GetGinTraceId(ginCtx)

				if dataJson {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.GormMiddleware]准备使用{gormRecordJson}保存数据：%s", data)
					}
					c.gormRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				} else {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.GormMiddleware]准备使用{gormRecordXml}保存数据：%s", data)
					}
					c.gormRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				}
			}
		}()
	}
}
