package main

import (

	// "backend/routes"
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
	PORT := os.Getenv("PORT")
	// Connect to database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the User model
	for _, model := range models.AllModels {
		err = db.AutoMigrate(model)
	}
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}
	// Create a new Gin router
	r := gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                 // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},  // Allow all methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Allow all headers
		ExposeHeaders:    []string{"Content-Length"},                                    // Expose specific headers
		AllowCredentials: true,                                                          // Allow credentials (e.g., cookies)
		MaxAge:           12 * time.Hour,                                                // Cache preflight requests for 12 hours
	}))
	// Define routes
	routes.RegisterStudentRoutes(r, db)
	routes.RegisterPageRoutes(r, db)

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	log.Printf("Server started on  http://localhost:%v", PORT)
	log.Printf("Backend API docs started on %v", "http://localhost:4200/swagger/index.html")
	if err := r.Run(":" + PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
