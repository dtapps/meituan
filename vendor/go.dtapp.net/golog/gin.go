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
	"io/ioutil"
	"net"
	"os"
	"runtime"
)

// GinClientFun *GinClient 驱动
type GinClientFun func() *GinClient

// GinClientJsonFun *GinClient 驱动
// jsonStatus bool json状态
type GinClientJsonFun func() (*GinClient, bool)

// GinClient 框架
type GinClient struct {
	gormClient *dorm.GormClient // 数据库驱动
	ipService  *goip.Client     // ip服务
	zapLog     *ZapLog          // 日志服务
	logDebug   bool             // 日志开关
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

// GinClientConfig 框架实例配置
type GinClientConfig struct {
	IpService     *goip.Client            // ip服务
	GormClientFun dorm.GormClientTableFun // 日志配置
	Debug         bool                    // 日志开关
	ZapLog        *ZapLog                 // 日志服务
	JsonStatus    bool                    // json状态
}

// NewGinClient 创建框架实例化
// client 数据库服务
// tableName 表名
// ipService ip服务
func NewGinClient(config *GinClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	c.config.jsonStatus = config.JsonStatus

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

		c.ipService = config.IpService

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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (c *GinClient) jsonUnmarshal(data string) (result interface{}) {
	_ = json.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (c *GinClient) Middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取
		data, _ := ioutil.ReadAll(ginCtx.Request.Body)

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

			var traceId = gotrace_id.GetGinTraceId(ginCtx)

			// 记录
			if dataJson {
				if c.logDebug {
					c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{gormRecordJson}保存数据：%s", data)
				}
				c.gormRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
			} else {
				if c.logDebug {
					c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{gormRecordXml}保存数据：%s", data)
				}
				c.gormRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
			}
		}()
	}
}
