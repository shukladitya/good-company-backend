// cmd/api/main.go
package main

import (
	"log"
	"os"
	"serveMovies/internal/handlers"
	"serveMovies/internal/middleware"
	"serveMovies/internal/models"
	"serveMovies/internal/router"
	"serveMovies/internal/services"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	// dsn := "host=localhost user=movie password=movie dbname=moviedb port=5432 sslmode=disable"

	db_url := os.Getenv("DB_EXTERNAL_URL")

	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{})
	//
	//
	// Initialize services
	authService := &services.AuthService{DB: db}
	emailService := services.NewEmailService()
	passwordResetService := &services.PasswordResetService{DB: db}

	//
	//
	// Initialize handlers
	// Note: we need to inject all services in the handlers here only(Dependency injection) else it will not be initialized when using the handlers and give nil pointer error.
	authHandler := &handlers.AuthHandler{
		AuthService: authService,
		EmailService: emailService,
	}
	passwordResetHandler := &handlers.PasswordResetHandler{
		PasswordResetService: passwordResetService,
		EmailService: emailService,
	}

	//
	//
	// Initialize middleware
	authMiddleware := &middleware.AuthMiddleware{
		AuthService: authService,
	}



	// Initialize and setup router
	r := router.NewRouter(authHandler, passwordResetHandler, authMiddleware)
	app := r.Setup()

	app.Run(":8080")
}