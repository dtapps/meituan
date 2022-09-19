package golog

import "time"

// 模型
type ginPostgresqlLogString struct {
	LogId              uint      `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`               //【记录】编号
	TraceId            string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`                //【系统】跟踪编号
	RequestTime        time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`              //【请求】时间
	RequestUri         string    `gorm:"comment:【请求】请求链接 域名+路径+参数" json:"request_uri,omitempty"`          //【请求】请求链接 域名+路径+参数
	RequestUrl         string    `gorm:"comment:【请求】请求链接 域名+路径" json:"request_url,omitempty"`             //【请求】请求链接 域名+路径
	RequestApi         string    `gorm:"index;comment:【请求】请求接口 路径" json:"request_api,omitempty"`          //【请求】请求接口 路径
	RequestMethod      string    `gorm:"index;comment:【请求】请求方式" json:"request_method,omitempty"`          //【请求】请求方式
	RequestProto       string    `gorm:"comment:【请求】请求协议" json:"request_proto,omitempty"`                 //【请求】请求协议
	RequestUa          string    `gorm:"comment:【请求】请求UA" json:"request_ua,omitempty"`                    //【请求】请求UA
	RequestReferer     string    `gorm:"comment:【请求】请求referer" json:"request_referer,omitempty"`          //【请求】请求referer
	RequestBody        string    `gorm:"comment:【请求】请求主体" json:"request_body,omitempty"`                  //【请求】请求主体
	RequestUrlQuery    string    `gorm:"comment:【请求】请求URL参数" json:"request_url_query,omitempty"`          //【请求】请求URL参数
	RequestIp          string    `gorm:"index;comment:【请求】请求客户端Ip" json:"request_ip,omitempty"`           //【请求】请求客户端Ip
	RequestIpCountry   string    `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_country,omitempty"`   //【请求】请求客户端城市
	RequestIpProvince  string    `gorm:"index;comment:【请求】请求客户端省份" json:"request_ip_province,omitempty"`  //【请求】请求客户端省份
	RequestIpCity      string    `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_city,omitempty"`      //【请求】请求客户端城市
	RequestIpIsp       string    `gorm:"index;comment:【请求】请求客户端运营商" json:"request_ip_isp,omitempty"`      //【请求】请求客户端运营商
	RequestIpLongitude float64   `gorm:"index;comment:【请求】请求客户端经度" json:"request_ip_longitude,omitempty"` //【请求】请求客户端经度
	RequestIpLatitude  float64   `gorm:"index;comment:【请求】请求客户端纬度" json:"request_ip_latitude,omitempty"`  //【请求】请求客户端纬度
	//RequestIpLocation `gorm:"index;comment:【请求】请求客户端位置" json:"request_ip_location,omitempty"`  //【请求】请求客户端位置
	RequestHeader  string    `gorm:"comment:【请求】请求头" json:"request_header,omitempty"`          //【请求】请求头
	ResponseTime   time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`      //【返回】时间
	ResponseCode   int       `gorm:"index;comment:【返回】状态码" json:"response_code,omitempty"`     //【返回】状态码
	ResponseMsg    string    `gorm:"comment:【返回】描述" json:"response_msg,omitempty"`             //【返回】描述
	ResponseData   string    `gorm:"comment:【返回】数据" json:"response_data,omitempty"`            //【返回】数据
	CostTime       int64     `gorm:"comment:【系统】花费时间" json:"cost_time,omitempty"`              //【系统】花费时间
	SystemHostName string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`  //【系统】主机名
	SystemInsideIp string    `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"` //【系统】内网ip
	SystemOs       string    `gorm:"index;comment:【系统】系统类型" json:"system_os,omitempty"`        //【系统】系统类型
	SystemArch     string    `gorm:"index;comment:【系统】系统架构" json:"system_arch,omitempty"`      //【系统】系统架构
	GoVersion      string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`             //【程序】Go版本
	SdkVersion     string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`           //【程序】Sdk版本
}
