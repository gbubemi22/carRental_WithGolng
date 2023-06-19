package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Verify and parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid token signing method")
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Token is valid, continue to the next handler
		c.Next()
	}
}
