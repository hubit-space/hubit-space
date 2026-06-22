package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func PanicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v\n", r)
				debug.PrintStack()

				var errMsg string
				switch v := r.(type) {
				case string:
					errMsg = v
				case error:
					errMsg = v.Error()
				default:
					errMsg = fmt.Sprintf("%v", v)
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error: " + errMsg,
				})
			}
		}()

		c.Next()
	}
}