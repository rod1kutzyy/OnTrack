package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/logger"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logger.Logger.WithField("stack", stack).Errorf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
