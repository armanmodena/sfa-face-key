package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "*")                                                            // Allow all origins (you can restrict this to specific domains)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")                             // Allowed methods
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-With") // Allowed headers
		c.Header("Access-Control-Allow-Credentials", "true")                                                    // Allow credentials (cookies, etc.)
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Requested-With")                           // Expose additional headers

		// If the request is an OPTIONS request, immediately return a 200 status code
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}
