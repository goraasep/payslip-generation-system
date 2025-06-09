package dto

type CreateAttendancePeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"` // e.g. "2025-06-01"
	EndDate   string `json:"end_date" binding:"required"`   // e.g. "2025-06-30"
}
