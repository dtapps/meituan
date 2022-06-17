package meituan

import (
	"go.dtapp.net/gojson"
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
	"gorm.io/datatypes"
)

// 记录日志
func (app *App) postgresqlLog(request gorequest.Response) {
	app.log.Record(golog.ApiPostgresqlLog{
		RequestTime:           golog.TimeString{Time: request.RequestTime},                      //【请求】时间
		RequestUri:            request.RequestUri,                                               //【请求】链接
		RequestUrl:            gorequest.UriParse(request.RequestUri).Url,                       //【请求】链接
		RequestApi:            gorequest.UriParse(request.RequestUri).Path,                      //【请求】接口
		RequestMethod:         request.RequestMethod,                                            //【请求】方式
		RequestParams:         datatypes.JSON(gojson.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
		RequestHeader:         datatypes.JSON(gojson.JsonEncodeNoError(request.RequestHeader)),  //【返回】头部
		ResponseHeader:        datatypes.JSON(gojson.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                                       //【返回】状态码
		ResponseBody:          request.ResponseBody,                                             //【返回】内容
		ResponseContentLength: request.ResponseContentLength,                                    //【返回】大小
		ResponseTime:          golog.TimeString{Time: request.ResponseTime},                     //【返回】时间
	})
}
