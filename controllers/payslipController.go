package controllers

import (
	"fmt"
	"net/http"

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
