package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/routes"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	statsviz.RegisterDefault()

	go func() {
		fmt.Println("Point your browser to http://localhost:7000/debug/statsviz/")
		log.Fatal(http.ListenAndServe(":7000", nil))
	}()
}
func main() {
	config.ConnectDatabase()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080") // http://localhost:8080
}
