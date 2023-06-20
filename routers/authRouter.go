package routes

import (
	controller "github.com/gbubemi22/go-rentals/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
     authGroup := incomingRoutes.Group("/api/v1/auth") 
	{
		authGroup.POST("/signup", controller.CreateUser())
		authGroup.POST("/agent/signup", controller.CreateAgent())
		authGroup.POST("/login", controller.Login())
	}

	
}
