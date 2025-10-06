package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
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

var DB *sql.DB

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Create a new employee with the provided information
// @Tags employee
// @Accept json
// @Produce json
// @Param employee body Employee true "Employee object that needs to be created"
// @Success 201 {object} Employee
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error creating employee"
// @Router /employee [post]
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
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
	query := `INSERT INTO m_employee (employee_code, prefix_name, first_name, last_name, nickname, email, phone_number, gender, birth_date, hire_date, department, position, employment_type) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	err = DB.QueryRow(query, "", employee.PrefixName, employee.FirstName, employee.LastName, "", "", "", 0, nil, nil, "", "", 0).Scan(&employee.ID)
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
// @Tags employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID (UUID)"
// @Success 200 {object} Employee
// @Failure 400 {string} string "Employee ID is required"
// @Failure 404 {string} string "Employee not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error retrieving employee"
// @Router /employee/{id} [get]
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get employee ID from URL path
	// Extract ID from path like /api/employee/123
	path := r.URL.Path
	employeeID := path[len("/api/employee/"):]

	if employeeID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	// Query employee from database
	query := `SELECT id, employee_code, prefix_name, first_name, last_name, nickname, 
				email, phone_number, gender, birth_date, hire_date, department, 
				position, employment_type, is_active, created_at, updated_at 
			  FROM m_employee WHERE id = $1`

	var employee Employee
	var birthDate, hireDate, createdAt, updatedAt sql.NullTime
	var employeeCode, nickname, email, phoneNumber, department, position sql.NullString
	var gender, employmentType sql.NullInt32

	err := DB.QueryRow(query, employeeID).Scan(
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
