package models

import "gorm.io/gorm"

type Payslip struct {
	gorm.Model

	PayrollID uint    `gorm:"not null;index" json:"payroll_id"`
	Payroll   Payroll `gorm:"foreignKey:PayrollID"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	BaseSalary         float64 `gorm:"not null;default:0"`
	ProratedSalary     float64 `gorm:"not null;default:0"`
	OvertimePay        float64 `gorm:"not null;default:0"`
	OvertimeCount      int     `gorm:"not null;default:0"`
	OvertimeHours      float64 `gorm:"not null;default:0"`
	AttendanceCount    int     `gorm:"not null;default:0"`
	AttendancePeriod   int     `gorm:"not null;default:0"`
	ReimbursementTotal float64 `gorm:"not null;default:0"`
	TakeHomePay        float64 `gorm:"not null;default:0"`

	Reimbursements []PayslipReimbursement `gorm:"foreignKey:PayslipID"`
}
