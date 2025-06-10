package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendancePeriod struct {
	gorm.Model

	StartDate time.Time `gorm:"not null;type:date" json:"start_date"` // store only the date
	EndDate   time.Time `gorm:"not null;type:date" json:"end_date"`   // store only the date

	// LockedAt *time.Time `gorm:"type:timestamp" json:"locked_at,omitempty"` // nullable; full timestamp
}
