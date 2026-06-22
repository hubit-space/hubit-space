package middleware

import (
	"hubit-space/service/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		SECURITY_CODE_HUBIT := config.GetEnv("SECURITY_CODE_HUBIT", "91b637d8fcd2c6da6359e6963113a1170de795e4b725b84d1e0b4ck19skdj8fm")

		securityCode := c.GetHeader("Security-Code")
		if securityCode == "" {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		if securityCode != SECURITY_CODE_HUBIT {
			c.JSON(401, gin.H{"error": "You are not authorized to access this service"})
			c.Abort()
			return
		}

		c.Next()
	}
}
