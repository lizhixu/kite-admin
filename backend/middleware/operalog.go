package middleware

import (
	"backend/config"
	"backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Stop timer
		endTime := time.Now()
		latency := endTime.Sub(startTime).Milliseconds()

		// Skip logging for GET syslog/list to avoid log spam when viewing logs
		if c.Request.URL.Path == "/api/syslog/list" {
			return
		}

		// Extract user from context
		userIdVal, exists := c.Get("userID")
		userId := uint(0)
		if exists {
			if id, ok := userIdVal.(uint); ok {
				userId = id
			} else if idFloat, ok := userIdVal.(float64); ok {
				userId = uint(idFloat)
			}
		}

		usernameVal, exists := c.Get("username")
		username := ""
		if exists && usernameVal != nil {
			username = usernameVal.(string)
		}

		userAgent := strings.Join(c.Request.Header["User-Agent"], " ")
		if len(userAgent) > 250 {
			userAgent = userAgent[:250]
		}

		logEntry := models.SysLog{
			UserID:     userId,
			Username:   username,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			IP:         c.ClientIP(),
			StatusCode: c.Writer.Status(),
			Latency:    latency,
			UserAgent:  userAgent,
		}

		// Fast async insert
		go func(l models.SysLog) {
			config.DB.Create(&l)
		}(logEntry)
	}
}
