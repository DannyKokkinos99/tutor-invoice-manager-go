package main

import (
	"log"
	"os"
	"time"
	"tutor-invoice-manager/config"
	_ "tutor-invoice-manager/docs"
	"tutor-invoice-manager/models"
	"tutor-invoice-manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Testing swagger API for GO
// @version 0.70
// @description Testing go for rest api using Gin
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4200
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Ensures logs go to the container logs
	log.SetOutput(os.Stdout)
	// Load the .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Variables
	projectType := os.Getenv("PROJECT_TYPE")
	domain := os.Getenv("DOMAIN")
	port := os.Getenv("PORT")

	if projectType == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := config.InitDB(projectType)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if projectType == "development" {
		// Auto-migrate the User model
		for _, model := range models.AllModels {
			err = db.AutoMigrate(model)
		}
		if err != nil {
			log.Fatalf("Failed to migrate schema: %v", err)
		}
	}
	// Create a new Gin router
	r := gin.Default()

	// Configure CORS middleware
	if projectType == "production" {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, //TODO: Change to the hosted url on cloud run
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	} else if projectType == "development" {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	// Define routes
	routes.RegisterStudentRoutes(r, db)
	routes.RegisterPageRoutes(r, db)
	routes.RegisterInvoiceRoutes(r, db)
	routes.RegisterDashboardRoutes(r, db)

	if projectType == "development" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Printf("Server started on %s:%s", domain, port)
		log.Printf("Backend API docs started on %s:%s/swagger/index.html", domain, port)
	}
	if projectType == "production" { //dev works using air
		log.Printf("Production server started on %s:%s", domain, port)
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
