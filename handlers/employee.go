package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

type EmployeeListResponse struct {
	Data       []Employee `json:"data"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
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

// GetEmployeeList godoc
// @Summary Get list of employees with pagination and sorting
// @Description Get paginated list of employees with optional sorting
// @Tags employee
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param sort_by query string false "Sort by field (id, first_name, last_name, created_at, etc.)"
// @Param sort_order query string false "Sort order (asc, desc) default: asc"
// @Param search query string false "Search in first_name, last_name, email"
// @Success 200 {object} EmployeeListResponse
// @Failure 400 {string} string "Invalid parameters"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error retrieving employees"
// @Router /employees [get]
func GetEmployeeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Pagination parameters
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	// Sorting parameters
	sortBy := query.Get("sort_by")
	sortOrder := query.Get("sort_order")
	search := query.Get("search")

	// Validate sort fields
	validSortFields := map[string]bool{
		"id": true, "employee_code": true, "prefix_name": true,
		"first_name": true, "last_name": true, "email": true,
		"department": true, "position": true, "created_at": true,
		"updated_at": true, "is_active": true,
	}

	if sortBy == "" {
		sortBy = "created_at" // Default sort
	}

	if !validSortFields[sortBy] {
		http.Error(w, "Invalid sort field", http.StatusBadRequest)
		return
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc" // Default order
	}

	// Build the WHERE clause for search
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = " WHERE (first_name ILIKE $" + strconv.Itoa(argIndex) +
			" OR last_name ILIKE $" + strconv.Itoa(argIndex+1) +
			" OR email ILIKE $" + strconv.Itoa(argIndex+2) + ")"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM m_employee" + whereClause
	var totalRecords int
	err := DB.QueryRow(countQuery, args...).Scan(&totalRecords)
	if err != nil {
		http.Error(w, "Error counting employees: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build the main query
	mainQuery := fmt.Sprintf(`
		SELECT id, employee_code, prefix_name, first_name, last_name, nickname, 
			   email, phone_number, gender, birth_date, hire_date, department, 
			   position, employment_type, is_active, created_at, updated_at 
		FROM m_employee%s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, sortBy, strings.ToUpper(sortOrder), argIndex, argIndex+1)

	// Add limit and offset to args
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := DB.Query(mainQuery, args...)
	if err != nil {
		http.Error(w, "Error querying employees: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse results
	var employees []Employee
	for rows.Next() {
		var employee Employee
		var birthDate, hireDate, createdAt, updatedAt sql.NullTime
		var employeeCode, nickname, email, phoneNumber, department, position sql.NullString
		var gender, employmentType sql.NullInt32

		err := rows.Scan(
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

		if err != nil {
			http.Error(w, "Error scanning employee: "+err.Error(), http.StatusInternalServerError)
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

		employees = append(employees, employee)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating employees: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalRecords + pageSize - 1) / pageSize

	// Build response
	response := EmployeeListResponse{
		Data:       employees,
		Total:      totalRecords,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Description Update employee details by employee ID
// @Tags employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID (UUID)"
// @Param employee body Employee true "Updated employee object"
// @Success 200 {object} Employee
// @Failure 400 {string} string "Invalid request body or Employee ID is required"
// @Failure 404 {string} string "Employee not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error updating employee"
// @Router /employee/{id} [put]
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get employee ID from URL path
	path := r.URL.Path
	employeeID := path[len("/api/employee/"):]

	if employeeID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
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

	// Update employee in database
	query := `UPDATE m_employee 
			  SET employee_code = $1, prefix_name = $2, first_name = $3, last_name = $4, 
			      nickname = $5, email = $6, phone_number = $7, gender = $8, 
			      birth_date = $9, hire_date = $10, department = $11, position = $12, 
			      employment_type = $13, is_active = $14, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $15
			  RETURNING id, employee_code, prefix_name, first_name, last_name, nickname, 
			            email, phone_number, gender, birth_date, hire_date, department, 
			            position, employment_type, is_active, created_at, updated_at`

	// Prepare nullable values
	var birthDate, hireDate interface{}
	if employee.BirthDate != "" {
		birthDate = employee.BirthDate
	}
	if employee.HireDate != "" {
		hireDate = employee.HireDate
	}

	var updatedEmployee Employee
	var updatedBirthDate, updatedHireDate, createdAt, updatedAt sql.NullTime
	var employeeCode, nickname, email, phoneNumber, department, position sql.NullString
	var gender, employmentType sql.NullInt32

	err = DB.QueryRow(query,
		employee.EmployeeCode,
		employee.PrefixName,
		employee.FirstName,
		employee.LastName,
		employee.Nickname,
		employee.Email,
		employee.PhoneNumber,
		employee.Gender,
		birthDate,
		hireDate,
		employee.Department,
		employee.Position,
		employee.EmploymentType,
		employee.IsActive,
		employeeID,
	).Scan(
		&updatedEmployee.ID,
		&employeeCode,
		&updatedEmployee.PrefixName,
		&updatedEmployee.FirstName,
		&updatedEmployee.LastName,
		&nickname,
		&email,
		&phoneNumber,
		&gender,
		&updatedBirthDate,
		&updatedHireDate,
		&department,
		&position,
		&employmentType,
		&updatedEmployee.IsActive,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error updating employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle nullable fields
	if employeeCode.Valid {
		updatedEmployee.EmployeeCode = employeeCode.String
	}
	if nickname.Valid {
		updatedEmployee.Nickname = nickname.String
	}
	if email.Valid {
		updatedEmployee.Email = email.String
	}
	if phoneNumber.Valid {
		updatedEmployee.PhoneNumber = phoneNumber.String
	}
	if gender.Valid {
		updatedEmployee.Gender = int(gender.Int32)
	}
	if updatedBirthDate.Valid {
		updatedEmployee.BirthDate = updatedBirthDate.Time.Format("2006-01-02")
	}
	if updatedHireDate.Valid {
		updatedEmployee.HireDate = updatedHireDate.Time.Format("2006-01-02")
	}
	if department.Valid {
		updatedEmployee.Department = department.String
	}
	if position.Valid {
		updatedEmployee.Position = position.String
	}
	if employmentType.Valid {
		updatedEmployee.EmploymentType = int(employmentType.Int32)
	}
	if createdAt.Valid {
		updatedEmployee.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
	}
	if updatedAt.Valid {
		updatedEmployee.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
	}

	// Return updated employee
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedEmployee)
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Delete employee by employee ID
// @Tags employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID (UUID)"
// @Success 200 {object} map[string]string "Employee deleted successfully"
// @Failure 400 {string} string "Employee ID is required"
// @Failure 404 {string} string "Employee not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Error deleting employee"
// @Router /employee/{id} [delete]
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get employee ID from URL path
	path := r.URL.Path
	employeeID := path[len("/api/employee/"):]

	if employeeID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	// Delete employee from database
	query := `DELETE FROM m_employee WHERE id = $1`

	result, err := DB.Exec(query, employeeID)
	if err != nil {
		http.Error(w, "Error deleting employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if employee was found and deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Return success response
	response := map[string]string{
		"message": "Employee deleted successfully",
		"id":      employeeID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
