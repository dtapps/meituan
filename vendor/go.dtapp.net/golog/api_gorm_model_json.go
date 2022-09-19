package golog

import (
	"gorm.io/datatypes"
	"time"
)

// 模型
type apiPostgresqlLogJson struct {
	LogId                 uint           `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`           //【记录】编号
	TraceId               string         `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`            //【系统】跟踪编号
	RequestTime           time.Time      `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`          //【请求】时间
	RequestUri            string         `gorm:"comment:【请求】链接" json:"request_uri,omitempty"`                 //【请求】链接
	RequestUrl            string         `gorm:"comment:【请求】链接" json:"request_url,omitempty"`                 //【请求】链接
	RequestApi            string         `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`           //【请求】接口
	RequestMethod         string         `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`        //【请求】方式
	RequestParams         datatypes.JSON `gorm:"type:jsonb;comment:【请求】参数" json:"request_params,omitempty"`   //【请求】参数
	RequestHeader         datatypes.JSON `gorm:"type:jsonb;comment:【请求】头部" json:"request_header,omitempty"`   //【请求】头部
	RequestIp             string         `gorm:"index;comment:【请求】请求Ip" json:"request_ip,omitempty"`          //【请求】请求Ip
	ResponseHeader        datatypes.JSON `gorm:"type:jsonb;comment:【返回】头部" json:"response_header,omitempty"`  //【返回】头部
	ResponseStatusCode    int            `gorm:"index;comment:【返回】状态码" json:"response_status_code,omitempty"` //【返回】状态码
	ResponseBody          datatypes.JSON `gorm:"type:jsonb;comment:【返回】数据" json:"response_content,omitempty"` //【返回】数据
	ResponseContentLength int64          `gorm:"comment:【返回】大小" json:"response_content_length,omitempty"`     //【返回】大小
	ResponseTime          time.Time      `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`         //【返回】时间
	SystemHostName        string         `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`     //【系统】主机名
	SystemInsideIp        string         `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"`    //【系统】内网ip
	SystemOs              string         `gorm:"index;comment:【系统】系统类型" json:"system_os,omitempty"`           //【系统】系统类型
	SystemArch            string         `gorm:"index;comment:【系统】系统架构" json:"system_arch,omitempty"`         //【系统】系统架构
	GoVersion             string         `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`                //【程序】Go版本
	SdkVersion            string         `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`              //【程序】Sdk版本
}
