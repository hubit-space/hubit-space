package middleware

import (
	"hubit-space/gateway/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtSecret := []byte(config.GetEnv("JWT_SECRET_KEY", ""))

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if ctx.Request.URL.Path != "/refresh-token" && (err != nil || !token.Valid) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if userId, ok := payload["user_id"].(string); ok {
			ctx.Request.Header.Set("UserId", userId)
		}
		if name, ok := payload["name"].(string); ok {
			ctx.Request.Header.Set("Name", name)
		}
		if email, ok := payload["email"].(string); ok {
			ctx.Request.Header.Set("Email", email)
		}
		if userCode, ok := payload["user_code"].(string); ok {
			ctx.Request.Header.Set("UserCode", userCode)
		}
		if roleCode, ok := payload["role_code"].(string); ok {
			ctx.Request.Header.Set("RoleCode", roleCode)
		}

		ctx.Next()
	}
}
