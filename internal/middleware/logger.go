package middleware

import (
	"log"
	"road2ca/pkg/minigin"
	"time"
)

type LoggerMiddleware interface {
	SettingLogger(c *minigin.Context)
}

type loggerMiddleware struct{}

func NewLoggerMiddleware() LoggerMiddleware {
	return &loggerMiddleware{}
}

// SettingLogger ログを出力するミドルウェア
func (m *loggerMiddleware) SettingLogger(c *minigin.Context) {
	startTime := time.Now()

	c.Next()
	duration := time.Since(startTime)
	statusCode := c.Writer.Status()
	log.Printf("Method: %s, Path: %s, Status: %d, Duration: %s",
		c.Request.Method, c.Request.URL.Path, statusCode, duration)
}
