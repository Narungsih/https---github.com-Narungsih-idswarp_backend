package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "backend/docs"

	_ "github.com/lib/pq"
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

type Employee struct {
	ID             string `json:"id"`
	EmployeeCode   string `json:"employee_code"`
	PrefixName     string `json:"prefix_name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Nickname       string `json:"nickname"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Gender         int    `json:"gender"`
	BirthDate      string `json:"birth_date"`
	HireDate       string `json:"hire_date"`
	Department     string `json:"department"`
	Position       string `json:"position"`
	EmploymentType int    `json:"employment_type"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=IDS-warp sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error verifying connection to database:", err)
	}

	// Create employees table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS m_employees (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		employee_code VARCHAR(20),
		prefix_name VARCHAR(50) NOT NULL,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		nickname VARCHAR(50),
		email VARCHAR(150),
		phone_number VARCHAR(50),
		gender SMALLINT DEFAULT 0,
		birth_date DATE,
		hire_date DATE,
		department VARCHAR(150),
		position VARCHAR(150),
		employment_type SMALLINT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("Database connection established and table created successfully")
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Create a new employee with the provided information
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body Employee true "Employee object that needs to be created"
// @Success 201 {object} Employee
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error creating employee"
// @Router /create/employees [post]
func createEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var employee Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if employee.PrefixName == "" || employee.FirstName == "" || employee.LastName == "" {
		http.Error(w, "prefix_name, first_name and last_name are required", http.StatusBadRequest)
		return
	}

	// Insert employee into database
	query := `INSERT INTO m_employees (employee_code, prefix_name, first_name, last_name, nickname, email, phone_number, gender, birth_date, hire_date, department, position, employment_type) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	err = db.QueryRow(query, "", employee.PrefixName, employee.FirstName, employee.LastName, "", "", "", 0, nil, nil, "", "", 0).Scan(&employee.ID)
	if err != nil {
		http.Error(w, "Error creating employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created employee
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Setup routes
	http.HandleFunc("/api/create/employees", enableCORS(createEmployee))

	// Swagger route
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
