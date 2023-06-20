package main

import (
	middleware "github.com/gbubemi22/go-rentals/middleware"
	routes "github.com/gbubemi22/go-rentals/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	//"strconv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())
     

	routes.UserRoutes(router)
	routes.AuthRoutes(router)
	routes.CarRoutes(router)
	router.NoRoute(middleware.NotFound())
	//router.Use(middleware.ErrorHandlerMiddleware())
     

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.Run(":" + port)
}
