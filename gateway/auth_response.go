package middleware

import (
	"hubit-space/gateway/config"

	"github.com/gin-gonic/gin"
)

func GeneralMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		PUBLIC_SECURITY_CODE_SCB_WEBSITE := config.GetEnv("PUBLIC_SECURITY_CODE_SCB_WEBSITE", "9FkdC7B15A946Jd26S191NwM0L8F0S3K0D74A6E7I95CB90D1E7O6M4S9G21L83963A1170De795E4B725B8B637D8FcD2C6XqPZrWk")

		securityCode := c.GetHeader("Security-Code")
		if securityCode == "" {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		if securityCode != PUBLIC_SECURITY_CODE_SCB_WEBSITE {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		c.Next()
	}
}
