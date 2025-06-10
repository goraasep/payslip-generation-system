package dto

type CreateAttendanceLogRequest struct {
	AttendancePeriodID uint   `json:"attendance_period_id" binding:"required"`
	Date               string `json:"date" binding:"required"`
}
