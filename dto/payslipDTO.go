package dto

type PayslipRequest struct {
	PayrollID uint `json:"payroll_id"`
}

type PayslipResponse struct {
	ID                 uint                   `json:"id"`
	PayrollID          uint                   `json:"payroll_id"`
	UserID             uint                   `json:"user_id"`
	BaseSalary         float64                `json:"base_salary"`
	ProratedSalary     float64                `json:"prorated_salary"`
	OvertimePay        float64                `json:"overtime_pay"`
	OvertimeCount      int                    `json:"overtime_count"`
	OvertimeHours      float64                `json:"overtime_hours"`
	ReimbursementTotal float64                `json:"reimbursement_total"`
	Reimbursements     []ReimburseLogResponse `json:"reimbursements"`
	AttendanceCount    int                    `json:"attendance_count"`
	AttendancePeriod   int                    `json:"attendance_period"` // e.g., number of working days
	TakeHomePay        float64                `json:"take_home_pay"`
	CreatedAt          string                 `json:"created_at"` // formatted datetime
}
