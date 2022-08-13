package gotrace_id

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gostring"
)

// CustomTraceIdContext 自定义设置跟踪编号上下文
func CustomTraceIdContext() context.Context {
	return context.WithValue(context.Background(), "trace_id", gostring.GetUuId())
}

// SetGinTraceIdContext 设置跟踪编号上下文
func SetGinTraceIdContext(c *gin.Context) context.Context {
	return context.WithValue(context.Background(), "trace_id", GetGinTraceId(c))
}

// GetTraceIdContext 通过上下文获取跟踪编号
func GetTraceIdContext(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value("trace_id"))
}
