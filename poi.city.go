package meituan

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"net/http"
)

type PoiCityResponse struct {
	Code int `json:"code"` // 状态码 0表示请求正常
	Data []struct {
		Pinyin string `json:"pinyin"` // 城市拼音
		Name   string `json:"name"`   // 城市名称
		ID     int    `json:"id"`     // 城市id
	} `json:"data"` // 返回城市列表
}

type PoiCityResult struct {
	Result PoiCityResponse    // 结果
	Body   []byte             // 内容
	Http   gorequest.Response // 请求
}

func newPoiCityResult(result PoiCityResponse, body []byte, http gorequest.Response) *PoiCityResult {
	return &PoiCityResult{Result: result, Body: body, Http: http}
}

// PoiCity 基础数据 - 开放城市接口
// https://openapi.meituan.com/#api-0.%E5%9F%BA%E7%A1%80%E6%95%B0%E6%8D%AE-GetHttpsOpenapiMeituanComPoiCity
func (c *Client) PoiCity(ctx context.Context, notMustParams ...gorequest.Params) (*PoiCityResult, error) {
	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	// 请求
	request, err := c.request(ctx, apiUrl+"/poi/city", params, http.MethodGet)
	if err != nil {
		return newPoiCityResult(PoiCityResponse{}, request.ResponseBody, request), err
	}
	// 定义
	var response PoiCityResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	return newPoiCityResult(response, request.ResponseBody, request), err
}
