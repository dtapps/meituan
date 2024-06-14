package meituan

import (
	"context"
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
}

func newApiMiniCodeResult(result ApiMiniCodeResponse, body []byte, http gorequest.Response) *ApiMiniCodeResult {
	return &ApiMiniCodeResult{Result: result, Body: body, Http: http}
}

// ApiMiniCode 小程序生成二维码（新版）
// https://union.meituan.com/v2/apiDetail?id=26
func (c *Client) ApiMiniCode(ctx context.Context, notMustParams ...gorequest.Params) (*ApiMiniCodeResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "api/miniCode")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("appkey", c.GetAppKey()) // 媒体名称，可在推广者备案-媒体管理中查询
	params.Set("sign", c.getSign(c.GetSecret(), params))

	// 请求
	var response ApiMiniCodeResponse
	request, err := c.request(ctx, "api/miniCode", params, http.MethodGet, &response)
	return newApiMiniCodeResult(response, request.ResponseBody, request), err
}
