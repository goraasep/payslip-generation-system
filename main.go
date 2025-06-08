package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/models"
	"github.com/goraasep/payslip-generation-system/routes"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})

	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":8080") // http://localhost:8080
}
