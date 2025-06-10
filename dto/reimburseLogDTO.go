package dto

type CreateReimburseLogRequest struct {
	AttendancePeriodID uint    `json:"attendance_period_id" binding:"required"`
	Date               string  `json:"date" binding:"required"`
	Description        string  `json:"description" binding:"required"`
	Amount             float64 `json:"amount" binding:"required"`
}

type ReimburseLogResponse struct {
	ID                 uint    `json:"id"`
	AttendancePeriodID uint    `json:"attendance_period_id"`
	UserID             uint    `json:"user_id"`
	Date               string  `json:"date"`
	Description        string  `json:"description"`
	Amount             float64 `json:"amount"`
}
