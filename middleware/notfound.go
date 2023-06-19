package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route does not exist"})
	}
}
