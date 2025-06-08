package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/models"
)

func AdminDashboard(c *gin.Context) {
	response.Success(c, "Welcome to the Admin Dashboard!", nil)
}

// Example: list all users
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Roles").Find(&users).Error; err != nil {
		response.BadRequest(c, "Failed to fetch users")
		return
	}

	response.Success(c, "All users retrieved successfully", users)
}
