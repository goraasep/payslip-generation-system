package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceLog struct {
	gorm.Model

	AttendancePeriodID uint             `json:"attendance_period_id"`
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"`

	UserID uint `gorm:"not null;" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Date time.Time `gorm:"not null;type:date;" json:"date"`
}
