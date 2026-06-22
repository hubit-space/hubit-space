package middleware

import (
	"hubit-space/gateway/config"

	"github.com/gin-gonic/gin"
)

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		PUBLIC_SECURITY_CODE_SCB := config.GetEnv("PUBLIC_SECURITY_CODE_SCB", "3fkdcb7d15a946jd26s191nwm0l8f0s3k0d74a6e7i95cb90d1e7o6m4s9g21l83963a1170de795e4b725b8b637d8fcd2c6")

		securityCode := c.GetHeader("Security-Code")
		if securityCode == "" {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		if securityCode != PUBLIC_SECURITY_CODE_SCB {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		c.Next()
	}
}
