package goip

import (
	"go.dtapp.net/goip/geoip"
	"go.dtapp.net/goip/ip2region"
	"go.dtapp.net/goip/ip2region_v2"
	"go.dtapp.net/goip/ipv6wry"
	"go.dtapp.net/goip/qqwry"
	"net"
	"strconv"
)

type AnalyseResult struct {
	Ip              string                   `json:"ip,omitempty"`
	QqwryInfo       qqwry.QueryResult        `json:"qqwry_info"`
	Ip2regionInfo   ip2region.QueryResult    `json:"ip2region_info"`
	Ip2regionV2info ip2region_v2.QueryResult `json:"ip2regionv2_info"`
	GeoipInfo       geoip.QueryCityResult    `json:"geoip_info"`
	Ipv6wryInfo     ipv6wry.QueryResult      `json:"ipv6wry_info"`
}

func (c *Client) Analyse(item string) AnalyseResult {
	isIp := c.isIpv4OrIpv6(item)
	ipByte := net.ParseIP(item)
	switch isIp {
	case ipv4:
		qqeryInfo, _ := c.QueryQqWry(ipByte)
		ip2regionInfo, _ := c.QueryIp2Region(ipByte)
		ip2regionV2Info, _ := c.QueryIp2RegionV2(ipByte)
		geoipInfo, _ := c.QueryGeoIp(ipByte)
		return AnalyseResult{
			Ip:              ipByte.String(),
			QqwryInfo:       qqeryInfo,
			Ip2regionInfo:   ip2regionInfo,
			Ip2regionV2info: ip2regionV2Info,
			GeoipInfo:       geoipInfo,
		}
	case ipv6:
		geoipInfo, _ := c.QueryGeoIp(ipByte)
		ipv6Info, _ := c.QueryIpv6wry(ipByte)
		return AnalyseResult{
			Ip:          ipByte.String(),
			GeoipInfo:   geoipInfo,
			Ipv6wryInfo: ipv6Info,
		}
	default:
		return AnalyseResult{}
	}
}

// CheckIpv4 检查数据是不是IPV4
func (c *Client) CheckIpv4(ips string) bool {
	if len(ips) > 3 {
		return false
	}
	nums, err := strconv.Atoi(ips)
	if err != nil {
		return false
	}
	if nums < 0 || nums > 255 {
		return false
	}
	if len(ips) > 1 && ips[0] == '0' {
		return false
	}
	return true
}

// CheckIpv6 检测是不是IPV6
func (c *Client) CheckIpv6(ips string) bool {
	if ips == "" {
		return true
	}
	if len(ips) > 4 {
		return false
	}
	for _, val := range ips {
		if !((val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F')) {
			return false
		}
	}
	return true
}
