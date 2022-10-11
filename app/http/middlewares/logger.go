package middlewares

import (
	"bytes"
	"io"
	"thub/pkg/helpers"
	"thub/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// 这个结构体的作用是继承了 gin.ResponseWriter结构的所有属性和方法 并额外增加一个字段 body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 自定义 Logger 中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 response 内容
		w := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个buffer对象 只能读取一次
			requestBody, _ := io.ReadAll(c.Request.Body)
			// 读取后重新赋值，以供后续操作
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 设置开始时间
		start := time.Now()
		c.Next()

		// 回调 开始记录日志的逻辑
		cost := time.Since(start) // time.Since == time.Now().Sub(start)
		responseStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", responseStatus),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", helpers.MicrosecondsStr(cost)),
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			// 请求的内容
			logFields = append(logFields, zap.String("Request Body", string(requestBody)))
			// 响应的内容
			logFields = append(logFields, zap.String("Response Body", w.body.String()))
		}

		if responseStatus > 400 && responseStatus <= 499 {
			logger.Warn("HTTP Waring"+cast.ToString(responseStatus), logFields...)
		} else if responseStatus >= 500 && responseStatus <= 499 {
			logger.Error("HTTP Error"+cast.ToString(responseStatus), logFields...)
		} else {
			logger.Debug("HTTP Access Log", logFields...)
		}
	}
}
