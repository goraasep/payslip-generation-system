package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goraasep/payslip-generation-system/models"

	"golang.org/x/crypto/bcrypt"
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

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to PostgreSQL database!")
	}

	DB = database
	DB.AutoMigrate(&models.User{}, &models.Role{}, &models.AttendancePeriod{}, &models.AttendanceLog{}, &models.ReimburseLog{}, &models.OvertimeLog{})
	Seeding()
}

func Seeding() {
	// Create roles
	adminRole := models.Role{Name: "ADMIN"}
	userRole := models.Role{Name: "USER"}

	DB.FirstOrCreate(&adminRole, models.Role{Name: "ADMIN"})
	DB.FirstOrCreate(&userRole, models.Role{Name: "USER"})

	// ===== Admin User =====
	hashedAdminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "admin",
		Email:    "admin@admin.com",
		Password: string(hashedAdminPassword),
		Roles:    []*models.Role{&adminRole},
	}
	DB.FirstOrCreate(&admin, models.User{Email: "admin@admin.com"})
	DB.Model(&admin).Association("Roles").Append(&adminRole)

	// ===== Normal User =====
	hashedUserPassword, _ := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "user",
		Email:    "user@user.com",
		Password: string(hashedUserPassword),
		Roles:    []*models.Role{&userRole},
	}
	DB.FirstOrCreate(&user, models.User{Email: "user@user.com"})
	DB.Model(&user).Association("Roles").Append(&userRole)
}
