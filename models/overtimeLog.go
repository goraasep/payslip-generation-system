package models

import (
	"time"

	"gorm.io/gorm"
)

type OvertimeLog struct {
	gorm.Model

	AttendancePeriodID uint             `json:"attendance_period_id"`
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Date        time.Time `gorm:"not null;type:date;index" json:"date"`               // ensures no time part, indexed for fast querying
	Description string    `gorm:"type:text" json:"description"`                       // allows long text if needed
	Hour        int       `gorm:"not null;check:hour >= 1 AND hour <= 3" json:"hour"` // business rule constraint
}
