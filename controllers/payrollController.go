package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/dto"
	"github.com/goraasep/payslip-generation-system/models"
)

func RunPayroll(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Check if user is admin
	var currentUser models.User
	if err := config.DB.Preload("Roles").First(&currentUser, userID).Error; err != nil {
		response.BadRequest(c, "User not found")
		return
	}
	isAdmin := false
	for _, role := range currentUser.Roles {
		if role.Name == "ADMIN" {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		response.Unauthorized(c, "Only admin can run payroll")
		return
	}

	var input dto.PayrollRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tx := config.DB.Begin() // Begin transaction

	// Find period
	var period models.AttendancePeriod
	if err := tx.First(&period, input.AttendancePeriodID).Error; err != nil {
		tx.Rollback()
		response.BadRequest(c, "Attendance period not found")
		return
	}

	// Check if payroll already processed
	var existing models.Payroll
	if err := tx.Where("attendance_period_id = ?", input.AttendancePeriodID).First(&existing).Error; err == nil {
		tx.Rollback()
		response.BadRequest(c, "Payroll already processed for this period")
		return
	}

	// Create payroll record
	payroll := models.Payroll{
		AttendancePeriodID: input.AttendancePeriodID,
		ProcessedAt:        time.Now(),
	}
	if err := tx.Create(&payroll).Error; err != nil {
		tx.Rollback()
		response.InternalError(c, "Failed to create payroll")
		return
	}

	// Fetch all users
	var users []models.User
	if err := tx.Preload("Roles").Find(&users).Error; err != nil {
		tx.Rollback()
		response.InternalError(c, "Failed to fetch users")
		return
	}

	for _, user := range users {
		// Skip admin
		isUserAdmin := false
		for _, role := range user.Roles {
			if role.Name == "ADMIN" {
				isUserAdmin = true
				break
			}
		}
		if isUserAdmin {
			continue
		}

		baseSalary := user.Salary

		// Count attendance days
		var presentCount int64
		tx.Model(&models.AttendanceLog{}).
			Where("user_id = ? AND attendance_period_id = ?", user.ID, input.AttendancePeriodID).
			Count(&presentCount)

		// Working days
		workDays := countWeekdays(period.StartDate, period.EndDate)
		if workDays == 0 {
			workDays = 1
		}
		proratedSalary := baseSalary * float64(presentCount) / float64(workDays)

		// Overtime hours
		var totalOvertime float64
		tx.Model(&models.OvertimeLog{}).
			Where("user_id = ? AND attendance_period_id = ?", user.ID, input.AttendancePeriodID).
			Select("COALESCE(SUM(hour), 0)").Scan(&totalOvertime)

		hourlyRate := baseSalary / float64(workDays*8)
		overtimePay := totalOvertime * hourlyRate * 2

		// Reimbursements
		var reimbursements []models.ReimburseLog
		tx.Where("user_id = ? AND attendance_period_id = ?", user.ID, input.AttendancePeriodID).
			Find(&reimbursements)

		var reimburseTotal float64
		for _, r := range reimbursements {
			reimburseTotal += r.Amount
		}

		takeHome := proratedSalary + overtimePay + reimburseTotal

		// Create payslip
		payslip := models.Payslip{
			PayrollID:          payroll.ID,
			UserID:             user.ID,
			BaseSalary:         baseSalary,
			ProratedSalary:     proratedSalary,
			OvertimePay:        overtimePay,
			ReimbursementTotal: reimburseTotal,
			TakeHomePay:        takeHome,
		}
		if err := tx.Create(&payslip).Error; err != nil {
			tx.Rollback()
			response.InternalError(c, fmt.Sprintf("Failed to create payslip for user %s", user.Email))
			return
		}

		for _, r := range reimbursements {
			payslipReimburse := models.PayslipReimbursement{
				PayslipID:      payslip.ID,
				Description:    r.Description,
				Amount:         r.Amount,
				ReimburseLogID: r.ID,
			}
			if err := tx.Create(&payslipReimburse).Error; err != nil {
				tx.Rollback()
				response.InternalError(c, "Failed to create payslip reimbursement")
				return
			}
		}
	}
	payrollResponse := dto.PayrollResponse{
		ID:                 payroll.ID,
		AttendancePeriodID: payroll.AttendancePeriodID,
		ProcessedAt:        payroll.ProcessedAt.Format("2006-01-02 15:04:05"),
	}

	tx.Commit()
	response.Success(c, "Payroll successfully processed", payrollResponse)
}

func countWeekdays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			count++
		}
	}
	return count
}
