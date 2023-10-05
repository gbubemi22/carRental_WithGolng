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
		carGroup.GET("/:car_id", controller.GetOneCar())
		carGroup.PATCH("/:car_id", controller.UpdateCar())
		carGroup.DELETE("/:id", controller.DeleteCar())
		carGroup.POST("/:car_id/images", controller.UpLoadCarImage())
		carGroup.GET("/find/:user_id", controller.GetAllCarsByAgent())

	}
}
