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

func GetAllOvertimeLogs(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	search, _ := c.GetQuery("search")
	order := c.DefaultQuery("order", "desc")
	field := c.DefaultQuery("field", "id")

	userID := c.MustGet("user_id").(uint)
	//check if user admin or not
	var user models.User
	err := config.DB.Preload("Roles").First(&user, userID).Error
	if err != nil {
		response.BadRequest(c, "User not found")
		return
	}
	isAdmin := false
	for _, role := range user.Roles {
		if role.Name == "ADMIN" {
			isAdmin = true
			break
		}

	}

	var overtimeLogs []models.OvertimeLog
	query := config.DB.Model(&models.OvertimeLog{})

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if search != "" {
		query = query.Where("description ILIKE ?", "%"+search+"%")
	}
	var total int64
	query.Count(&total)

	query = query.Order(fmt.Sprintf("%s %s", field, order))

	err = query.Offset(start).Limit(length).Find(&overtimeLogs).Error
	if err != nil {
		response.BadRequest(c, "Failed to fetch overtime logs")
		return
	}

	paginationResponse := dto.PaginationResponse{
		RecordsFiltered: total,
		Data:            overtimeLogs,
	}

	response.Success(c, "Overtime logs retrieved", paginationResponse)
}

func CreateOvertimeLog(c *gin.Context) {
	var input dto.CreateOvertimeLogRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uint)

	layout := "2006-01-02" // Format: YYYY-MM-DD
	date, err := time.Parse(layout, input.Date)
	if err != nil {
		response.BadRequest(c, "Invalid date format. Use YYYY-MM-DD.")
		return
	}

	var attendance models.AttendanceLog
	err = config.DB.Where("attendance_period_id = ? AND user_id = ? AND date = ?", input.AttendancePeriodID, userID, date).First(&attendance).Error
	if err != nil {
		response.BadRequest(c, "Attendance log not found, please create attendance log first")
		return
	}

	var existing models.OvertimeLog
	err = config.DB.Where("attendance_period_id = ? AND user_id = ? AND date = ?", input.AttendancePeriodID, userID, date).First(&existing).Error
	if err == nil {
		response.BadRequest(c, "Overtime log already exists for this date")
		return
	}

	var period models.AttendancePeriod
	if err := config.DB.First(&period, input.AttendancePeriodID).Error; err != nil {
		response.BadRequest(c, "Attendance period not found")
		return
	}
	if date.Before(period.StartDate) || date.After(period.EndDate) {
		response.BadRequest(c, "Date is outside the attendance period range")
		return
	}

	overtimeLog := models.OvertimeLog{
		Description:        input.Description,
		AttendancePeriodID: input.AttendancePeriodID,
		Date:               date,
		UserID:             userID,
		Hour:               input.Hour,
	}

	if err := config.DB.Create(&overtimeLog).Error; err != nil {
		response.InternalError(c, "Failed to create overtime log")
		return
	}

	response.Success(c, "Overtime log created", overtimeLog)
}
