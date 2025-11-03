package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		logEntry := logger.Logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency.Milliseconds(),
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"query":       query,
			"user_agent":  userAgent,
		})

		if statusCode >= 500 {
			logEntry.Error("Server error")
		} else if statusCode >= 400 {
			logEntry.Warn("Client error")
		} else {
			logEntry.Info("Request processed")
		}
	}
}
