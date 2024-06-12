package meituan

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

type ApiGenerateLinkResponse struct {
	Status int    `json:"status"`         // 状态值，0为成功，非0为异常
	Des    string `json:"des,omitempty"`  // 异常描述信息
	Data   string `json:"data,omitempty"` // 最终的推广链接
}

type ApiGenerateLinkResult struct {
	Result ApiGenerateLinkResponse // 结果
	Body   []byte                  // 内容
	Http   gorequest.Response      // 请求
}

func newApiGenerateLinkResult(result ApiGenerateLinkResponse, body []byte, http gorequest.Response) *ApiGenerateLinkResult {
	return &ApiGenerateLinkResult{Result: result, Body: body, Http: http}
}

// ApiGenerateLink 自助取链接口（新版）
// https://union.meituan.com/v2/apiDetail?id=25
func (c *Client) ApiGenerateLink(ctx context.Context, notMustParams ...gorequest.Params) (*ApiGenerateLinkResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "api/generateLink")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("appkey", c.GetAppKey()) // 媒体名称，可在推广者备案-媒体管理中查询
	params.Set("sign", c.getSign(c.GetSecret(), params))

	// 请求
	request, err := c.request(ctx, "api/generateLink", params, http.MethodGet)
	if err != nil {
		if c.trace {
			c.span.SetStatus(codes.Error, err.Error())
		}
		return newApiGenerateLinkResult(ApiGenerateLinkResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response ApiGenerateLinkResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil && c.trace {
		c.span.SetStatus(codes.Error, err.Error())
	}
	return newApiGenerateLinkResult(response, request.ResponseBody, request), err
}
