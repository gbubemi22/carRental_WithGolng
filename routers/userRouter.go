package routes

import (
	controller "github.com/gbubemi22/go-rentals/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userGroup := incomingRoutes.Group("/api/v1/users")
	{
		userGroup.GET("/", controller.GetAllUsers())
		userGroup.GET("/:user_id", controller.GetUser())
	}
}
