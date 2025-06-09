package models

import (
	"time"

	"gorm.io/gorm"
)

type ReimburseLog struct {
	gorm.Model
	AttendancePeriodID uint             `json:"attendance_period_id"`
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"` // relation

	UserID      uint      `json:"user_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}
