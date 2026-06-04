package middleware

import (
	"time"

	"blog/server/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)
		l := logger.L
		if l == nil {
			return
		}

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Duration("duration", duration),
			zap.String("request_id", GetRequestID(c)),
		}
		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		l.Info("request", fields...)
	}
}
