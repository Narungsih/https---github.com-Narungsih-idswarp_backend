package main

import (
	"log"
	"net/http"
	"os"

	_ "backend/docs"

	"backend/database"
	"backend/handlers"
	"backend/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Employee Management API
// @version 1.0
// @description API for managing employees with PostgreSQL database
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.Close()

	// Share database connection with handlers
	handlers.DB = database.DB

	// Setup routes
	http.HandleFunc("/api/create/employees", middleware.EnableCORS(handlers.CreateEmployee))
	http.HandleFunc("/api/employees", middleware.EnableCORS(handlers.GetEmployeeByID))

	// Swagger route
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	serverAddr := ":" + port
	log.Printf("Server starting on port %s", serverAddr)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
