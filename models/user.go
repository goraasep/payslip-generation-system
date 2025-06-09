package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:text;not null" json:"-"`

	Salary float64 `gorm:"not null;default:0" json:"salary"` // monthly base salary

	Roles []*Role `gorm:"many2many:user_roles;" json:"roles"`
}
