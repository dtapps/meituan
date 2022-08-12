package meituan

import (
	"context"
	"encoding/json"
	"go.dtapp.net/gorequest"
	"net/http"
)

type ApiMiniCodeResponse struct {
	Status int    `json:"status"`         // 状态值，0为成功，非0为异常
	Des    string `json:"des,omitempty"`  // 异常描述信息
	Data   string `json:"data,omitempty"` // 小程序二维码图片地址
}

type ApiMiniCodeResult struct {
	Result ApiMiniCodeResponse // 结果
	Body   []byte              // 内容
	Http   gorequest.Response  // 请求
	Err    error               // 错误
}

func newApiMiniCodeResult(result ApiMiniCodeResponse, body []byte, http gorequest.Response, err error) *ApiMiniCodeResult {
	return &ApiMiniCodeResult{Result: result, Body: body, Http: http, Err: err}
}

// ApiMiniCode 小程序生成二维码（新版）
// https://union.meituan.com/v2/apiDetail?id=26
func (c *Client) ApiMiniCode(ctx context.Context, actId int64, sid string) *ApiMiniCodeResult {
	// 参数
	param := gorequest.NewParams()
	param.Set("appkey", c.config.AppKey)
	param.Set("sid", sid)
	param.Set("actId", actId)
	// 转换
	params := gorequest.NewParamsWith(param)
	params["sign"] = c.getSign(c.config.Secret, params)
	// 请求
	request, err := c.request(ctx, apiUrl+"/api/miniCode", params, http.MethodGet)
	// 定义
	var response ApiMiniCodeResponse
	err = json.Unmarshal(request.ResponseBody, &response)
	return newApiMiniCodeResult(response, request.ResponseBody, request, err)
}
