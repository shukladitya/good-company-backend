// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"
	"theGoodCompany/internal/handlers"
	"theGoodCompany/internal/middleware"
	"theGoodCompany/internal/models"
	"theGoodCompany/internal/router"
	"theGoodCompany/internal/services"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Mongo DB
	ctx := context.TODO()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        panic(err)
    }
    defer client.Disconnect(ctx)


	// Auto migrate models
	db.AutoMigrate(&models.User{})
	//Mongo DB
	collection := client.Database("yourDB").Collection("yourCollection")
	//
	//
	// Initialize services
	authService := &services.AuthService{DB: db}
	emailService := services.NewEmailService()
	passwordResetService := &services.PasswordResetService{DB: db}
	docService := services.NewDocumentService(collection)

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
	docHandler := &handlers.DocumentHandler{
		Service: docService,
	}

	//
	//
	// Initialize middleware
	authMiddleware := &middleware.AuthMiddleware{
		AuthService: authService,
	}



	// Initialize and setup router
	r := router.NewRouter(authHandler, passwordResetHandler, authMiddleware, docHandler)
	app := r.Setup()

	app.Run(":8080")
}