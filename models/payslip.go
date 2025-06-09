package models

import "gorm.io/gorm"

type Payslip struct {
	gorm.Model

	PayrollID uint    `gorm:"not null;index" json:"payroll_id"`
	Payroll   Payroll `gorm:"foreignKey:PayrollID" json:"payroll"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	BaseSalary         float64 `gorm:"not null" json:"base_salary"`         // full-month salary
	ProratedSalary     float64 `gorm:"not null" json:"prorated_salary"`     // based on attendance
	OvertimePay        float64 `gorm:"not null" json:"overtime_pay"`        // 2x salary rate * hours
	ReimbursementTotal float64 `gorm:"not null" json:"reimbursement_total"` // from PayslipReimbursements
	TakeHomePay        float64 `gorm:"not null" json:"take_home_pay"`       // total

	Reimbursements []PayslipReimbursement `gorm:"foreignKey:PayslipID" json:"reimbursements"`
}
