package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nanwannaporn/insurance-system/internal/api/handlers"
	"github.com/nanwannaporn/insurance-system/internal/domain"
	"github.com/nanwannaporn/insurance-system/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initializeGormDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL with GORM: %v", err)
	}

	err = db.AutoMigrate(
		&domain.Customer{},
		&domain.Beneficiaries{},
		&domain.HealthDeclaration{},
		&domain.Insurance{},
		&domain.CustomerInsurance{},
	)

	if err != nil {
		log.Fatalf("Failed to auto migrate database schema: %v", err)
	}
	log.Println("Database migration completed.")
	return db
}

func setupRoutes(r *gin.Engine, h *handlers.CustomerHandler) {
	api := r.Group("/api")
	{
		api.POST("/customer", h.CreateCustomerHandler)
		api.PUT("/customer/:id/beneficiary", h.UpdateBeneficiaryHandler)
		api.PUT("/customer/:id/health", h.UpdateHealthDeclarationHandler)
		api.GET("/plans", h.GetPlansHandler)
		api.POST("/policy/purchase", h.CreatePolicyHandler)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: Error loading .env file, using default environment variables.")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if username == "" {
		username = "myuser"
	}
	if password == "" {
		password = "mypassword"
	}
	if databaseName == "" {
		databaseName = "mydatabase"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, port, username, password, databaseName)

	gormDB := initializeGormDB(dsn)
	customerServiceInstance := service.NewCustomerService(gormDB)
	customerHandler := handlers.NewCustomerHandler(customerServiceInstance)

	//route handler
	router := gin.Default()
	setupRoutes(router, customerHandler)
	router.Run() // listens on 0.0.0.0:8080
}
