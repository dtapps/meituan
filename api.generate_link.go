package meituan

import (
	"context"
	"encoding/json"
	"go.dtapp.net/gorequest"
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
	Err    error                   // 错误
}

func newApiGenerateLinkResult(result ApiGenerateLinkResponse, body []byte, http gorequest.Response, err error) *ApiGenerateLinkResult {
	return &ApiGenerateLinkResult{Result: result, Body: body, Http: http, Err: err}
}

// ApiGenerateLink 自助取链接口（新版）
// https://union.meituan.com/v2/apiDetail?id=25
func (c *Client) ApiGenerateLink(ctx context.Context, actId int64, sid string, linkType, shortLink int) *ApiGenerateLinkResult {
	// 参数
	param := gorequest.NewParams()
	param.Set("actId", actId)            // 活动id，可以在联盟活动列表中查看获取
	param.Set("appkey", c.config.AppKey) // 媒体名称，可在推广者备案-媒体管理中查询
	param.Set("sid", sid)                // 推广位sid，支持通过接口自定义创建，不受平台200个上限限制，长度不能超过64个字符，支持小写字母和数字，历史已创建的推广位不受这个约束
	param.Set("linkType", linkType)      // 投放链接的类型
	param.Set("shortLink", shortLink)    // 获取长链还是短链
	// 转换
	params := gorequest.NewParamsWith(param)
	params["sign"] = c.getSign(c.config.Secret, params)
	// 请求
	request, err := c.request(ctx, apiUrl+"/api/generateLink", params, http.MethodGet)
	// 定义
	var response ApiGenerateLinkResponse
	err = json.Unmarshal(request.ResponseBody, &response)
	return newApiGenerateLinkResult(response, request.ResponseBody, request, err)
}
