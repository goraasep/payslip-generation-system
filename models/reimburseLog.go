package models

import (
	"time"

	"gorm.io/gorm"
)

type ReimburseLog struct {
	gorm.Model

	AttendancePeriodID uint             `gorm:"not null" json:"attendance_period_id"`
	AttendancePeriod   AttendancePeriod `gorm:"foreignKey:AttendancePeriodID" json:"attendance_period"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Date        time.Time `gorm:"not null;type:date;index" json:"date"`     // store only the date
	Description string    `gorm:"type:text" json:"description"`             // allows long text
	Amount      float64   `gorm:"not null;check:amount >= 0" json:"amount"` // ensures non-negative value
}
