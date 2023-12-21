package middlewares

import "github.com/gin-gonic/gin"

func API() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
