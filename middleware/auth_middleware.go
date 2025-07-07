package middleware

import (
	"arkan-face-key/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		securityCode := config.SECURITY_CODE
		authHeader := c.GetHeader("Security-Code")

		if authHeader == "" {
			c.JSON(401, gin.H{
				"status": 401,
				"error":  "Security-Code header is required",
			})
			c.Abort()
			return
		}

		if authHeader != securityCode {
			c.JSON(401, gin.H{
				"status": 401,
				"error":  "Invalid Security-Code",
			})
			c.Abort()
			return
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}
