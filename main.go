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

// @title Employee Management & Location API
// @version 1.0
// @description API for managing employees with bilingual support and location dropdown data
// @host localhost:8080
// @BasePath /api

// employeeHandler routes requests based on HTTP method
func employeeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's a single employee operation (has ID in path)
	path := r.URL.Path
	if len(path) > len("/api/employee/") && path != "/api/employee" {
		// Operations on specific employee by ID
		switch r.Method {
		case http.MethodGet:
			handlers.GetEmployeeByID(w, r)
		case http.MethodPut:
			handlers.UpdateEmployee(w, r)
		case http.MethodDelete:
			handlers.DeleteEmployee(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		// Operations on employee collection
		switch r.Method {
		case http.MethodPost:
			handlers.CreateEmployee(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.Close()

	// Share database connection with handlers
	handlers.DB = database.DB

	// Setup routes
	http.HandleFunc("/api/employee", middleware.EnableCORS(employeeHandler))
	http.HandleFunc("/api/employee/", middleware.EnableCORS(employeeHandler))
	http.HandleFunc("/api/employees", middleware.EnableCORS(handlers.GetEmployeeList))

	// Location dropdown routes
	http.HandleFunc("/api/geographies", middleware.EnableCORS(handlers.GetGeographies))
	http.HandleFunc("/api/provinces", middleware.EnableCORS(handlers.GetProvinces))
	http.HandleFunc("/api/districts", middleware.EnableCORS(handlers.GetDistricts))
	http.HandleFunc("/api/subdistricts", middleware.EnableCORS(handlers.GetSubDistricts))

	// Department routes
	http.HandleFunc("/api/departments", middleware.EnableCORS(handlers.GetDepartments))
	http.HandleFunc("/api/positions", middleware.EnableCORS(handlers.GetPositions))

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
