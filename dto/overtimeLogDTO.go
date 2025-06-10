package dto

type CreateOvertimeLogRequest struct {
	AttendancePeriodID uint   `json:"attendance_period_id" binding:"required"`
	Date               string `json:"date" binding:"required"`
	Hour               int    `json:"hour" binding:"required"`
	Description        string `json:"description" binding:"required"`
}
