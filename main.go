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

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.Close()

	// Share database connection with handlers
	handlers.DB = database.DB

	// Setup routes
	http.HandleFunc("/api/employee", middleware.EnableCORS(handlers.CreateEmployee))
	http.HandleFunc("/api/employee/", middleware.EnableCORS(handlers.GetEmployeeByID))

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
