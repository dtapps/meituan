package goip

import "strings"

var (
	ipv4 = "IPV4"
	ipv6 = "IPV6"
)

func (c *Client) isIpv4OrIpv6(ip string) string {
	if len(ip) < 7 {
		return ""
	}
	arrIpv4 := strings.Split(ip, ".")
	if len(arrIpv4) == 4 {
		//. 判断IPv4
		for _, val := range arrIpv4 {
			if !c.CheckIpv4(val) {
				return ""
			}
		}
		return ipv4
	}
	arrIpv6 := strings.Split(ip, ":")
	if len(arrIpv6) == 8 {
		// 判断Ipv6
		for _, val := range arrIpv6 {
			if !c.CheckIpv6(val) {
				return "Neither"
			}
		}
		return ipv6
	}
	return ""
}
