package dto

type PayslipSummaryItem struct {
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	TakeHomePay float64 `json:"take_home_pay"`
}

type PayslipSummaryResponse struct {
	Payslips         []PayslipSummaryItem `json:"payslips"`
	TotalTakeHome    float64              `json:"total_take_home"`
	AttendancePeriod string               `json:"attendance_period"` // e.g., "2025-06-01 to 2025-06-30"
	ProcessedAt      string               `json:"processed_at"`      // formatted date string
}

type PayslipSummaryRequest struct {
	PayrollID uint `json:"payroll_id" binding:"required"`
}
