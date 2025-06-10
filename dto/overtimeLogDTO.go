package dto

type CreateOvertimeLogRequest struct {
	AttendancePeriodID uint   `json:"attendance_period_id" binding:"required"`
	Date               string `json:"date" binding:"required"`
	Hour               int    `json:"hour" binding:"required"`
	Description        string `json:"description" binding:"required"`
}

type OvertimeLogResponse struct {
	ID                 uint   `json:"id"`
	AttendancePeriodID uint   `json:"attendance_period_id"`
	UserID             uint   `json:"user_id"`
	Date               string `json:"date"`
	Description        string `json:"description"`
	Hour               int    `json:"hour"`
}
