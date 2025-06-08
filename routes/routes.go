package routes

import (
	"github.com/goraasep/payslip-generation-system/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", controllers.GetUsers)
		userRoutes.POST("/", controllers.CreateUser)
	}
}
