package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/goraasep/payslip-generation-system/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}
	DB = db

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.AttendancePeriod{},
		&models.AttendanceLog{},
		&models.ReimburseLog{},
		&models.OvertimeLog{},
		&models.Payroll{},
		&models.Payslip{},
		&models.PayslipReimbursement{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func Seeding() {
	adminRole := models.Role{Name: "ADMIN"}
	userRole := models.Role{Name: "USER"}
	DB.FirstOrCreate(&adminRole, models.Role{Name: "ADMIN"})
	DB.FirstOrCreate(&userRole, models.Role{Name: "USER"})

	hashedAdminPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}
	admin := models.User{
		Name:     "admin",
		Email:    "admin@admin.com",
		Password: string(hashedAdminPassword),
		Roles:    []*models.Role{&adminRole},
		Salary:   0,
	}
	DB.FirstOrCreate(&admin, models.User{Email: admin.Email})
	DB.Model(&admin).Association("Roles").Replace(&adminRole)

	// 3) 100 fake users
	hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash user password: %v", err)
	}
	const minSalary, maxSalary = 3_000_000, 10_000_000

	for i := 1; i <= 100; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		name := fmt.Sprintf("user%d", i)

		raw := rand.Intn(maxSalary-minSalary+1) + minSalary

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashedUserPassword),
			Roles:    []*models.Role{&userRole},
			Salary:   float64(raw),
		}

		DB.FirstOrCreate(&user, models.User{Email: email})
		DB.Model(&user).Association("Roles").Replace(&userRole)
	}

	log.Println("Seeding: roles + 1 admin + 100 users created (if not already present).")
}
