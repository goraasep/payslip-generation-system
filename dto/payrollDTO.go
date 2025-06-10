package dto

type PayrollRequest struct {
	AttendancePeriodID uint `json:"attendance_period_id"`
}

type PayrollResponse struct {
	ID                 uint   `json:"id"`
	AttendancePeriodID uint   `json:"attendance_period_id"`
	ProcessedAt        string `json:"processed_at"`
}
