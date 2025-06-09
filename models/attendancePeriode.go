package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendancePeriod struct {
	gorm.Model
	StartDate time.Time  `json:"start_date"`          // expects RFC3339 format
	EndDate   time.Time  `json:"end_date"`            // expects RFC3339 format
	LockedAt  *time.Time `json:"locked_at,omitempty"` // nullable
}
