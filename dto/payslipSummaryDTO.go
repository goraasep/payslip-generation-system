package dto

type PayslipSummaryItem struct {
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	TakeHomePay float64 `json:"take_home_pay"`
}

type PayslipSummaryResponse struct {
	Payslips      []PayslipSummaryItem `json:"payslips"`
	TotalTakeHome float64              `json:"total_take_home"`
}

type PayslipSummaryRequest struct {
	PayrollID uint `json:"payroll_id" binding:"required"`
}
