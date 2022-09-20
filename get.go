package meituan

import "go.dtapp.net/golog"

func (c *Client) GetAppKey() string {
	return c.config.appKey
}

func (c *Client) GetSecret() string {
	return c.config.secret
}

func (c *Client) GetLog() *golog.ApiClient {
	return c.log.client
}
