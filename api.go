package meituan

import (
	"context"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
)

func (c *Client) Request(ctx context.Context, url string, method string, notMustParams ...gorequest.Params) ([]byte, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, url)
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)

	// 请求
	request, err := c.request(ctx, url, params, method)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		c.TraceRecordError(err)
		return request.ResponseBody, err
	}

	return request.ResponseBody, err
}
