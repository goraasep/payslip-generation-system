package dto

type CreateAttendanceLogRequest struct {
	AttendancePeriodID uint   `json:"attendance_period_id" binding:"required"`
	Date               string `json:"date" binding:"required"`
}

type AttendanceLogResponse struct {
	ID                 uint   `json:"id"`
	AttendancePeriodID uint   `json:"attendance_period_id"`
	UserID             uint   `json:"user_id"`
	Date               string `json:"date"`
}
