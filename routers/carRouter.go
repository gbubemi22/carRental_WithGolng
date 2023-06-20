package routes

import (
	controller "github.com/gbubemi22/go-rentals/controllers"
	"github.com/gin-gonic/gin"
)

func CarRoutes(incomingRoutes *gin.Engine) {
	carGroup := incomingRoutes.Group("/api/v1")
	{
		carGroup.POST("/", controller.CreateCar())
	}
}
