package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"type:varchar(255);unique" json:"email"`
	Password string `gorm:"type:text" json:"password"`
}
