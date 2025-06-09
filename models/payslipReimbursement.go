package models

import "gorm.io/gorm"

type PayslipReimbursement struct {
	gorm.Model

	PayslipID uint    `gorm:"not null;index" json:"payroll_id"`
	Payslip   Payslip `gorm:"foreignKey:PayslipID" json:"-"`

	ReimburseLogID uint         `gorm:"not null;uniqueIndex" json:"reimburse_log_id"`
	ReimburseLog   ReimburseLog `gorm:"foreignKey:ReimburseLogID" json:"reimburse_log"`

	Description string  `gorm:"type:text" json:"description"`
	Amount      float64 `gorm:"not null;check:amount >= 0" json:"amount"`
}
