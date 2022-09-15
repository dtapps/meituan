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
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"time"
)

// GinGormClientConfig 框架实例配置
type GinGormClientConfig struct {
	IpService     *goip.Client     // ip服务
	GormClientFun ginGormClientFun // 日志配置
	Debug         bool             // 日志开关
	ZapLog        *ZapLog          // 日志服务
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

	c.log.gorm = true

	return c, nil
}

// 创建模型
func (c *GinClient) gormAutoMigrate() (err error) {
	err = c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLog{})
	if err != nil {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.gormAutoMigrate]：%s", err)
	}
	return err
}

// 模型结构体
type ginPostgresqlLog struct {
	LogId             uint      `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`              //【记录】编号
	TraceId           string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`               //【系统】跟踪编号
	RequestTime       time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`             //【请求】时间
	RequestUri        string    `gorm:"comment:【请求】请求链接 域名+路径+参数" json:"request_uri,omitempty"`         //【请求】请求链接 域名+路径+参数
	RequestUrl        string    `gorm:"comment:【请求】请求链接 域名+路径" json:"request_url,omitempty"`            //【请求】请求链接 域名+路径
	RequestApi        string    `gorm:"index;comment:【请求】请求接口 路径" json:"request_api,omitempty"`         //【请求】请求接口 路径
	RequestMethod     string    `gorm:"index;comment:【请求】请求方式" json:"request_method,omitempty"`         //【请求】请求方式
	RequestProto      string    `gorm:"comment:【请求】请求协议" json:"request_proto,omitempty"`                //【请求】请求协议
	RequestUa         string    `gorm:"comment:【请求】请求UA" json:"request_ua,omitempty"`                   //【请求】请求UA
	RequestReferer    string    `gorm:"comment:【请求】请求referer" json:"request_referer,omitempty"`         //【请求】请求referer
	RequestBody       string    `gorm:"comment:【请求】请求主体" json:"request_body,omitempty"`                 //【请求】请求主体
	RequestUrlQuery   string    `gorm:"comment:【请求】请求URL参数" json:"request_url_query,omitempty"`         //【请求】请求URL参数
	RequestIp         string    `gorm:"index;comment:【请求】请求客户端Ip" json:"request_ip,omitempty"`          //【请求】请求客户端Ip
	RequestIpCountry  string    `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_country,omitempty"`  //【请求】请求客户端城市
	RequestIpRegion   string    `gorm:"index;comment:【请求】请求客户端区域" json:"request_ip_region,omitempty"`   //【请求】请求客户端区域
	RequestIpProvince string    `gorm:"index;comment:【请求】请求客户端省份" json:"request_ip_province,omitempty"` //【请求】请求客户端省份
	RequestIpCity     string    `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_city,omitempty"`     //【请求】请求客户端城市
	RequestIpIsp      string    `gorm:"index;comment:【请求】请求客户端运营商" json:"request_ip_isp,omitempty"`     //【请求】请求客户端运营商
	RequestHeader     string    `gorm:"comment:【请求】请求头" json:"request_header,omitempty"`                //【请求】请求头
	ResponseTime      time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`            //【返回】时间
	ResponseCode      int       `gorm:"index;comment:【返回】状态码" json:"response_code,omitempty"`           //【返回】状态码
	ResponseMsg       string    `gorm:"comment:【返回】描述" json:"response_msg,omitempty"`                   //【返回】描述
	ResponseData      string    `gorm:"comment:【返回】数据" json:"response_data,omitempty"`                  //【返回】数据
	CostTime          int64     `gorm:"comment:【系统】花费时间" json:"cost_time,omitempty"`                    //【系统】花费时间
	SystemHostName    string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`        //【系统】主机名
	SystemInsideIp    string    `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"`       //【系统】内网ip
	SystemOs          string    `gorm:"index;comment:【系统】系统类型" json:"system_os,omitempty"`              //【系统】系统类型
	SystemArch        string    `gorm:"index;comment:【系统】系统架构" json:"system_arch,omitempty"`            //【系统】系统架构
	SystemCpuQuantity int       `gorm:"index;comment:【系统】CPU核数" json:"system_cpu_quantity,omitempty"`   //【系统】CPU核数
	GoVersion         string    `gorm:"index;comment:【程序】Go版本" json:"go_version,omitempty"`             //【程序】Go版本
	SdkVersion        string    `gorm:"index;comment:【程序】Sdk版本" json:"sdk_version,omitempty"`           //【程序】Sdk版本
}

// gormRecord 记录日志
func (c *GinClient) gormRecord(postgresqlLog ginPostgresqlLog) (err error) {

	postgresqlLog.SystemHostName = c.gormConfig.hostname
	postgresqlLog.SystemInsideIp = c.gormConfig.insideIp
	postgresqlLog.GoVersion = c.gormConfig.goVersion

	postgresqlLog.SdkVersion = Version

	postgresqlLog.SystemOs = c.config.os
	postgresqlLog.SystemArch = c.config.arch
	postgresqlLog.SystemCpuQuantity = c.config.maxProCs

	err = c.gormClient.Db.Table(c.gormConfig.tableName).Create(&postgresqlLog).Error
	if err != nil {
		c.zapLog.WithTraceIdStr(postgresqlLog.TraceId).Sugar().Errorf("[golog.gin.gormRecord]：%s", err)
	}

	return
}

func (c *GinClient) gormRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.gormRecordJson]收到保存数据要求：%s", c.gormConfig.tableName)
	}

	data := ginPostgresqlLog{
		TraceId:           traceId,                                                      //【系统】跟踪编号
		RequestTime:       requestTime,                                                  //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   dorm.JsonEncodeNoError(ginCtx.Request.URL.Query()),           //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                        //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     dorm.JsonEncodeNoError(ginCtx.Request.Header),                //【请求】请求头
		ResponseTime:      gotime.Current().Time,                                        //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      responseBody,                                                 //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
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

func (c *GinClient) gormRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.gormRecordXml]收到保存数据要求：%s", c.gormConfig.tableName)
	}

	data := ginPostgresqlLog{
		TraceId:           traceId,                                                      //【系统】跟踪编号
		RequestTime:       requestTime,                                                  //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   dorm.JsonEncodeNoError(ginCtx.Request.URL.Query()),           //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                        //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     dorm.JsonEncodeNoError(ginCtx.Request.Header),                //【请求】请求头
		ResponseTime:      gotime.Current().Time,                                        //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      responseBody,                                                 //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
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

// GormQuery 查询
func (c *GinClient) GormQuery(ctx context.Context) *gorm.DB {
	return c.gormClient.Db.Table(c.gormConfig.tableName)
}

// GormDelete 删除
func (c *GinClient) GormDelete(ctx context.Context, hour int64) error {
	return c.gormClient.Db.Table(c.gormConfig.tableName).Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&ginPostgresqlLog{}).Error
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
			if c.gormClient != nil && c.gormClient.Db != nil {

				var traceId = gotrace_id.GetGinTraceId(ginCtx)

				if dataJson {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.GormMiddleware]准备使用{gormRecordJson}保存数据：%s", data)
					}
					c.gormRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				} else {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.GormMiddleware]准备使用{gormRecordXml}保存数据：%s", data)
					}
					c.gormRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				}
			}
		}()
	}
}
