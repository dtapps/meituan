package goip

import (
	"go.dtapp.net/goip/geoip"
	"go.dtapp.net/goip/ip2region"
	"go.dtapp.net/goip/ip2region_v2"
	"go.dtapp.net/goip/ipv6wry"
	"go.dtapp.net/goip/qqwry"
)

type Client struct {
	ip2regionV2Client *ip2region_v2.Client
	ip2regionClient   *ip2region.Client
	qqwryClient       *qqwry.Client
	geoIpClient       *geoip.Client
	ipv6wryClient     *ipv6wry.Client
}

// NewIp 实例化
func NewIp() *Client {

	c := &Client{}

	c.ip2regionV2Client, _ = ip2region_v2.New()

	c.ip2regionClient = ip2region.New()

	c.qqwryClient = qqwry.New()

	c.geoIpClient, _ = geoip.New()

	c.ipv6wryClient = ipv6wry.New()

	return c
}

func (c *Client) Close() {
	c.geoIpClient.Close()
}
