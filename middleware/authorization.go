package middleware

import (
	"github.com/gin-gonic/gin"
	//"github.com/gbubemi22/go-rentals/models"
	"net/http"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user type"})
			c.Abort()
			return
		}

		authorized := false
		for _, allowedRole := range roles {
			if allowedRole == userTypeStr {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
