package routes

import (
	"github.com/goraasep/payslip-generation-system/controllers"

	"github.com/gin-gonic/gin"

	"github.com/goraasep/payslip-generation-system/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/me", controllers.Me)
}
