package routes

import (
	controller "github.com/gbubemi22/go-rentals/controllers"
	"github.com/gin-gonic/gin"
)

func CarRoutes(incomingRoutes *gin.Engine) {
	carGroup := incomingRoutes.Group("/api/v1/cars")
	{
		carGroup.POST("/", controller.CreateCar())
		carGroup.GET("", controller.GetAllCars())
		carGroup.GET("/:id", controller.GetOneCar())
		carGroup.PATCH("/:car_id", controller.UpdateCar())
		carGroup.DELETE("/:id", controller.DeleteCar())
		carGroup.POST("/:id/images", controller.UpdateCarImage())
		carGroup.GET("/find/:agentId", controller.GetAllCarsByAgent())

	}
}
