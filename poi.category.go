package meituan

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

type PoiCategoryResponse struct {
	Code int `json:"code"`
	Data []struct {
		Name    string `json:"name"`
		Subcate []struct {
			Name string `json:"name"` // 品类名称
			ID   int    `json:"id"`   // 品类id
		} `json:"subcate"`
		ID int `json:"id"`
	} `json:"data"`
}

type PoiCategoryResult struct {
	Result PoiCategoryResponse // 结果
	Body   []byte              // 内容
	Http   gorequest.Response  // 请求
}

func newPoiCategoryResult(result PoiCategoryResponse, body []byte, http gorequest.Response) *PoiCategoryResult {
	return &PoiCategoryResult{Result: result, Body: body, Http: http}
}

// PoiCategory 基础数据 - 品类接口
// https://openapi.meituan.com/#api-0.%E5%9F%BA%E7%A1%80%E6%95%B0%E6%8D%AE-GetHttpsOpenapiMeituanComPoiDistrictCityid1
func (c *Client) PoiCategory(ctx context.Context, cityID int, notMustParams ...gorequest.Params) (*PoiCategoryResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "poi/category")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("cityid", cityID)

	// 请求
	request, err := c.request(ctx, "poi/category", params, http.MethodGet)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		c.TraceRecordError(err)
		return newPoiCategoryResult(PoiCategoryResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response PoiCategoryResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		c.TraceRecordError(err)
	}
	return newPoiCategoryResult(response, request.ResponseBody, request), err
}
