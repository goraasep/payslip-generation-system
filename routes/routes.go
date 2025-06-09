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

	// All authenticated users
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", controllers.Me)
	}

	// Admin only
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.RequireRoles("ADMIN"))
	{
		adminGroup.GET("/users", controllers.GetAllUsers)
		adminGroup.GET("/attendance-periods", controllers.GetAllAttendancePeriods)
		adminGroup.POST("/attendance-periods", controllers.CreateAttendancePeriod)
	}

	// Admin and User
	authGroup := r.Group("/profile")
	authGroup.Use(middleware.AuthMiddleware(), middleware.RequireRoles("ADMIN", "USER"))
	{
		authGroup.GET("/me", controllers.Me)
	}
}
