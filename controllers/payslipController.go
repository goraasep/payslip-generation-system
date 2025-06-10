package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/dto"
	"github.com/goraasep/payslip-generation-system/models"
	"github.com/jung-kurt/gofpdf"
)

func GetPayslip(c *gin.Context) {
	var input dto.PayslipRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uint)

	var payslip models.Payslip
	err := config.DB.Preload("Reimbursements").
		Preload("Reimbursements.ReimburseLog").
		Preload("Payroll.AttendancePeriod").
		Preload("User").
		Where("payroll_id = ? AND user_id = ?", input.PayrollID, userID).
		First(&payslip).Error

	if err != nil {
		response.BadRequest(c, "Payslip not found")
		return
	}

	reimbursements := make([]dto.ReimburseLogResponse, 0)

	for _, r := range payslip.Reimbursements {
		reimbursements = append(reimbursements, dto.ReimburseLogResponse{
			ID:          r.ID,
			Description: r.Description,
			Amount:      r.Amount,
			Date:        r.ReimburseLog.Date.Format("2006-01-02"),
		})
	}

	payslipResponse := dto.PayslipResponse{
		ID:                 payslip.ID,
		PayrollID:          payslip.PayrollID,
		UserID:             payslip.UserID,
		BaseSalary:         payslip.BaseSalary,
		ProratedSalary:     payslip.ProratedSalary,
		OvertimePay:        payslip.OvertimePay,
		OvertimeCount:      payslip.OvertimeCount,
		OvertimeHours:      payslip.OvertimeHours,
		ReimbursementTotal: payslip.ReimbursementTotal,
		Reimbursements:     reimbursements,
		AttendanceCount:    payslip.AttendanceCount,
		AttendancePeriod:   payslip.AttendancePeriod,
		TakeHomePay:        payslip.TakeHomePay,
		CreatedAt:          payslip.CreatedAt.Format("2006-01-02"),
	}

	// Optional: if query param `pdf=true`, generate PDF
	if c.Query("pdf") == "true" {
		generatePayslipPDF(c, payslipResponse)
		return
	}

	response.Success(c, "Payslip retrieved", payslipResponse)
}

func generatePayslipPDF(c *gin.Context, data dto.PayslipResponse) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Payslip")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("User ID: %d", data.UserID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Payroll ID: %d", data.PayrollID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Created At: %s", data.CreatedAt))

	pdf.Ln(12)
	pdf.Cell(60, 10, fmt.Sprintf("Base Salary: Rp%.2f", data.BaseSalary))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Prorated Salary: Rp%.2f", data.ProratedSalary))

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 10, "Attendance:")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Attendance: %d / %d days", data.AttendanceCount, data.AttendancePeriod))

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 10, "Overtime:")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Overtime Count: %d times", data.OvertimeCount))
	pdf.Ln(6)
	pdf.Cell(60, 10, fmt.Sprintf("Overtime Hours: %.2f hours", data.OvertimeHours))
	pdf.Ln(6)
	pdf.Cell(60, 10, fmt.Sprintf("Overtime Pay: Rp%.2f", data.OvertimePay))

	pdf.Ln(10)
	pdf.Cell(60, 10, fmt.Sprintf("Reimbursement Total: Rp%.2f", data.ReimbursementTotal))
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Reimbursements:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	if len(data.Reimbursements) == 0 {
		pdf.Cell(60, 10, "No reimbursements.")
		pdf.Ln(8)
	} else {
		for _, r := range data.Reimbursements {
			pdf.Cell(60, 8, fmt.Sprintf("- %s (%s): Rp%.2f", r.Description, r.Date, r.Amount))
			pdf.Ln(6)
		}
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(60, 10, fmt.Sprintf("Take Home Pay: Rp%.2f", data.TakeHomePay))

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=payslip.pdf")
	err := pdf.Output(c.Writer)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate PDF"})
	}
}

func GetPayslipSummary(c *gin.Context) {
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
		response.Unauthorized(c, "Only admin can access this")
		return
	}

	// Get all payslips with user info
	var payslips []models.Payslip
	if err := config.DB.Preload("User").Find(&payslips).Error; err != nil {
		response.InternalError(c, "Failed to fetch payslips")
		return
	}

	// Build response
	var summary []dto.PayslipSummaryItem
	var totalTakeHome float64

	for _, payslip := range payslips {
		item := dto.PayslipSummaryItem{
			UserID:      payslip.UserID,
			UserName:    payslip.User.Name,
			TakeHomePay: payslip.TakeHomePay,
		}
		summary = append(summary, item)
		totalTakeHome += payslip.TakeHomePay
	}

	res := dto.PayslipSummaryResponse{
		Payslips:      summary,
		TotalTakeHome: totalTakeHome,
	}

	response.Success(c, "Payslip summary generated", res)
}

func GeneratePayslipSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Admin check
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
		response.Unauthorized(c, "Only admin can access this")
		return
	}

	// Bind body
	var req dto.PayslipSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Fetch payslips for the given payroll_id
	var payslips []models.Payslip
	if err := config.DB.Preload("User").
		Where("payroll_id = ?", req.PayrollID).
		Find(&payslips).Error; err != nil {
		response.InternalError(c, "Failed to fetch payslips")
		return
	}

	if len(payslips) == 0 {
		response.BadRequest(c, "No payslips found for the given payroll ID")
		return
	}

	var payroll models.Payroll
	if err := config.DB.Preload("AttendancePeriod").First(&payroll, req.PayrollID).Error; err != nil {
		response.BadRequest(c, "Payroll not found")
		return
	}

	start := payroll.AttendancePeriod.StartDate.Format("2006-01-02")
	end := payroll.AttendancePeriod.EndDate.Format("2006-01-02")
	processedAt := payroll.ProcessedAt.Format("2006-01-02 15:04:05")

	var total float64
	var summary []dto.PayslipSummaryItem
	for _, p := range payslips {
		total += p.TakeHomePay
		summary = append(summary, dto.PayslipSummaryItem{
			UserID:      p.UserID,
			UserName:    p.User.Name,
			TakeHomePay: p.TakeHomePay,
		})
	}

	// Optional: generate PDF
	if c.Query("pdf") == "true" {
		generatePayslipSummaryPDF(c, summary, total, fmt.Sprintf("%s to %s", start, end), processedAt)
		return
	}

	response.Success(c, "Payslip summary fetched successfully", dto.PayslipSummaryResponse{
		Payslips:         summary,
		TotalTakeHome:    total,
		AttendancePeriod: fmt.Sprintf("%s to %s", start, end),
		ProcessedAt:      processedAt,
	})

}

func generatePayslipSummaryPDF(c *gin.Context, summary []dto.PayslipSummaryItem, total float64, attendancePeriod string, processedAt string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Payslip Summary Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	// Display attendance period and processed date
	pdf.Cell(40, 10, fmt.Sprintf("Attendance Period: %s", attendancePeriod))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Processed At: %s", processedAt))
	pdf.Ln(8)

	// Display generated time
	now := time.Now().Format("2006-01-02 15:04:05")
	pdf.Cell(40, 10, fmt.Sprintf("Generated At: %s", now))
	pdf.Ln(12)

	// Table headers
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(60, 10, "Employee", "1", 0, "", false, 0, "")
	pdf.CellFormat(60, 10, "User ID", "1", 0, "", false, 0, "")
	pdf.CellFormat(60, 10, "Take Home Pay", "1", 1, "", false, 0, "")

	// Table content
	pdf.SetFont("Arial", "", 12)
	for _, s := range summary {
		pdf.CellFormat(60, 10, s.UserName, "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 10, fmt.Sprintf("%v", s.UserID), "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 10, fmt.Sprintf("Rp %.2f", s.TakeHomePay), "1", 1, "", false, 0, "")
	}

	// Total row
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(120, 10, "Total", "1", 0, "R", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("Rp %.2f", total), "1", 1, "", false, 0, "")

	// Output
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		response.InternalError(c, "Failed to generate PDF")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=payslip_summary.pdf")
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
