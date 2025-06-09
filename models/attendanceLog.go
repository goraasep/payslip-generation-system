package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceLog struct {
	gorm.Model
	AttendancePeriodID uint             `json:"attendance_period_id"`
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"` // relation

	UserID uint      `json:"user_id"`
	Date   time.Time `json:"date"`
}
