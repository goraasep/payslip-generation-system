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
		protected.GET("/attendance-periods", controllers.GetAllAttendancePeriods)
		protected.GET("/attendance-logs", controllers.GetAllAttendanceLogs)
		protected.GET("/overtime-logs", controllers.GetAllOvertimeLogs)
		protected.GET("/reimburse-logs", controllers.GetAllReimburseLogs)

		// Admin and User
		authGroup := protected.Group("/profile")
		authGroup.Use(middleware.RequireRoles("ADMIN", "USER"))
		{
			authGroup.GET("/me", controllers.Me)
		}

		// Admin only
		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RequireRoles("ADMIN"))
		{
			adminGroup.GET("/users", controllers.GetAllUsers)
			adminGroup.POST("/attendance-periods", controllers.CreateAttendancePeriod)
		}

		userGroup := protected.Group("/user")
		userGroup.Use(middleware.RequireRoles("USER"))
		{
			userGroup.POST("/attendance-logs", controllers.CreateAttendanceLog)
			userGroup.POST("/overtime-logs", controllers.CreateOvertimeLog)
			userGroup.POST("/reimburse-logs", controllers.CreateReimburseLog)
		}
	}
}
