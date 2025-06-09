package models

import (
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	gorm.Model

	AttendancePeriodID uint             `gorm:"not null;unique" json:"attendance_period_id"` // one payroll per period
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"`

	ProcessedAt time.Time `gorm:"not null" json:"processed_at"`
}
