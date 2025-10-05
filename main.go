package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "backend/docs"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

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
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	// Build connection string from environment variables
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

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

// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Get employee details by employee ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id query string true "Employee ID (UUID)"
// @Success 200 {object} Employee
// @Failure 400 {string} string "Employee ID is required"
// @Failure 404 {string} string "Employee not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error retrieving employee"
// @Router /employees [get]
func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get employee ID from query parameter
	employeeID := r.URL.Query().Get("id")
	if employeeID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	// Query employee from database
	query := `SELECT id, employee_code, prefix_name, first_name, last_name, nickname, 
				email, phone_number, gender, birth_date, hire_date, department, 
				position, employment_type, is_active, created_at, updated_at 
			  FROM m_employees WHERE id = $1`

	var employee Employee
	var birthDate, hireDate, createdAt, updatedAt sql.NullTime
	var employeeCode, nickname, email, phoneNumber, department, position sql.NullString
	var gender, employmentType sql.NullInt32

	err := db.QueryRow(query, employeeID).Scan(
		&employee.ID,
		&employeeCode,
		&employee.PrefixName,
		&employee.FirstName,
		&employee.LastName,
		&nickname,
		&email,
		&phoneNumber,
		&gender,
		&birthDate,
		&hireDate,
		&department,
		&position,
		&employmentType,
		&employee.IsActive,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error retrieving employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle nullable fields
	if employeeCode.Valid {
		employee.EmployeeCode = employeeCode.String
	}
	if nickname.Valid {
		employee.Nickname = nickname.String
	}
	if email.Valid {
		employee.Email = email.String
	}
	if phoneNumber.Valid {
		employee.PhoneNumber = phoneNumber.String
	}
	if gender.Valid {
		employee.Gender = int(gender.Int32)
	}
	if birthDate.Valid {
		employee.BirthDate = birthDate.Time.Format("2006-01-02")
	}
	if hireDate.Valid {
		employee.HireDate = hireDate.Time.Format("2006-01-02")
	}
	if department.Valid {
		employee.Department = department.String
	}
	if position.Valid {
		employee.Position = position.String
	}
	if employmentType.Valid {
		employee.EmploymentType = int(employmentType.Int32)
	}
	if createdAt.Valid {
		employee.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
	}
	if updatedAt.Valid {
		employee.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
	}

	// Return employee
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	http.HandleFunc("/api/employees", enableCORS(getEmployeeByID))

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
