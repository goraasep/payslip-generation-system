package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/dto"
	"github.com/goraasep/payslip-generation-system/models"
)

func AdminDashboard(c *gin.Context) {
	response.Success(c, "Welcome to the Admin Dashboard!", nil)
}

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

func GetAllAttendancePeriods(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	order := c.DefaultQuery("order", "desc")
	field := c.DefaultQuery("field", "id")

	var attendancePeriods []models.AttendancePeriod
	query := config.DB.Model(&models.AttendancePeriod{})

	var total int64
	query.Count(&total)

	query = query.Order(fmt.Sprintf("%s %s", field, order))

	err := query.Offset(start).Limit(length).Find(&attendancePeriods).Error
	if err != nil {
		response.BadRequest(c, "Failed to fetch attendance periods")
		return
	}

	paginationResponse := dto.PaginationResponse{
		RecordsFiltered: total,
		Data:            attendancePeriods,
	}

	response.Success(c, "Attendance periods retrieved", paginationResponse)
}

func CreateAttendancePeriod(c *gin.Context) {
	var input dto.CreateAttendancePeriodRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	layout := "2006-01-02" // the layout for parsing YYYY-MM-DD

	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD.")
		return
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD.")
		return
	}

	endDate = time.Date(
		endDate.Year(), endDate.Month(), endDate.Day(),
		23, 59, 59, 0, time.UTC,
	)

	attendancePeriod := models.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := config.DB.Create(&attendancePeriod).Error; err != nil {
		response.InternalError(c, "Failed to create attendance period")
		return
	}

	response.Success(c, "Attendance period created", attendancePeriod)
}
