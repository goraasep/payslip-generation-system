package controllers

import (
	"fmt"
	"strconv"

	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/dto"
	"github.com/goraasep/payslip-generation-system/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	search, _ := c.GetQuery("search")
	order := c.DefaultQuery("order", "desc")
	field := c.DefaultQuery("field", "id")
	var users []models.User
	query := config.DB.Model(&models.User{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	query.Count(&total)

	query = query.Order(fmt.Sprintf("%s %s", field, order))

	err := query.Offset(start).Limit(length).Find(&users).Error
	if err != nil {
		response.BadRequest(c, "Failed to fetch users")
		return
	}

	if err := config.DB.Preload("Roles").Find(&users).Error; err != nil {
		response.BadRequest(c, "Failed to fetch users")
		return
	}

	paginationResponse := dto.PaginationResponse{
		RecordsFiltered: total,
		Data:            users,
	}

	response.Success(c, "All users retrieved successfully", paginationResponse)
}
