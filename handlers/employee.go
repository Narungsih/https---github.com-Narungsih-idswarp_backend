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
	EmployeeID       string `json:"employee_id"`
	EmploymentType   int    `json:"employment_type"`
	Title            int    `json:"title"`
	FirstNameEN      string `json:"first_name_en"`
	LastNameEN       string `json:"last_name_en"`
	FirstNameTH      string `json:"first_name_th"`
	LastNameTH       string `json:"last_name_th"`
	NickNameEN       string `json:"nick_name_en"`
	NickNameTH       string `json:"nick_name_th"`
	PhoneNumber      string `json:"phone_number"`
	CompanyEmail     string `json:"company_email"`
	Nationality      string `json:"nationality"`
	Gender           int    `json:"gender"`
	TaxID            string `json:"tax_id"`
	BirthDate        string `json:"birth_date"`
	StartWorkDate    string `json:"start_work_date"`
	Status           int    `json:"status"`
	Remark           string `json:"remark"`
	Department       string `json:"department"`
	Position         string `json:"position"`
	Photo            string `json:"photo"`
	CustomAttributes string `json:"custom_attributes"`
	CreatedBy        string `json:"created_by"`
	CreatedDate      string `json:"created_date"`
	UpdatedBy        string `json:"updated_by,omitempty"`
	UpdatedDate      string `json:"updated_date,omitempty"`
	IsActive         bool   `json:"is_active"`
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
// @Description Create a new employee with bilingual information
// @Tags employee
// @Accept json
// @Produce json
// @Param employee body Employee true "Employee object"
// @Success 201 {object} Employee
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Server error"
// @Router /employee [post]
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if employee.CreatedBy == "" {
		employee.CreatedBy = "00000000-0000-0000-0000-000000000000"
	}

	query := `INSERT INTO m_employee (
		employment_type, title, first_name_en, last_name_en, first_name_th, last_name_th,
		nick_name_en, nick_name_th, phone_number, company_email, nationality, gender,
		tax_id, birth_date, start_work_date, status, remark, department, position,
		photo, custom_attributes, created_by, is_active
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::timestamp, $15::timestamp, $16, $17, $18, $19, $20, $21, $22::uuid, $23) 
	RETURNING employee_id, created_date`

	var createdDate sql.NullTime
	err := DB.QueryRow(query,
		employee.EmploymentType, employee.Title, employee.FirstNameEN, employee.LastNameEN,
		employee.FirstNameTH, employee.LastNameTH, employee.NickNameEN, employee.NickNameTH,
		employee.PhoneNumber, employee.CompanyEmail, employee.Nationality, employee.Gender,
		employee.TaxID, employee.BirthDate, employee.StartWorkDate, employee.Status,
		employee.Remark, employee.Department, employee.Position, employee.Photo,
		employee.CustomAttributes, employee.CreatedBy, employee.IsActive,
	).Scan(&employee.EmployeeID, &createdDate)

	if err != nil {
		http.Error(w, "Error creating employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if createdDate.Valid {
		employee.CreatedDate = createdDate.Time.Format("2023-01-15T00:00:00Z")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Get employee details by employee ID
// @Tags employee
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} Employee
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Server error"
// @Router /employee/{id} [get]
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employeeID := r.URL.Path[len("/api/employee/"):]
	if employeeID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT employee_id, employment_type, title, first_name_en, last_name_en, first_name_th, last_name_th,
		nick_name_en, nick_name_th, phone_number, company_email, nationality, gender, tax_id, birth_date, 
		start_work_date, status, remark, department, position, photo, custom_attributes, created_by, 
		created_date, updated_by, updated_date, is_active FROM m_employee WHERE employee_id = $1`

	var employee Employee
	var birthDate, startWorkDate, createdDate, updatedDate sql.NullTime
	var updatedBy sql.NullString

	err := DB.QueryRow(query, employeeID).Scan(
		&employee.EmployeeID, &employee.EmploymentType, &employee.Title,
		&employee.FirstNameEN, &employee.LastNameEN, &employee.FirstNameTH, &employee.LastNameTH,
		&employee.NickNameEN, &employee.NickNameTH, &employee.PhoneNumber, &employee.CompanyEmail,
		&employee.Nationality, &employee.Gender, &employee.TaxID, &birthDate, &startWorkDate,
		&employee.Status, &employee.Remark, &employee.Department, &employee.Position,
		&employee.Photo, &employee.CustomAttributes, &employee.CreatedBy, &createdDate,
		&updatedBy, &updatedDate, &employee.IsActive,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if birthDate.Valid {
		employee.BirthDate = birthDate.Time.Format("2006-01-02 15:04:05")
	}
	if startWorkDate.Valid {
		employee.StartWorkDate = startWorkDate.Time.Format("2006-01-02 15:04:05")
	}
	if createdDate.Valid {
		employee.CreatedDate = createdDate.Time.Format("2006-01-02 15:04:05")
	}
	if updatedBy.Valid {
		employee.UpdatedBy = updatedBy.String
	}
	if updatedDate.Valid {
		employee.UpdatedDate = updatedDate.Time.Format("2006-01-02 15:04:05")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

// GetEmployeeList godoc
// @Summary Get list of employees
// @Description Get paginated list of employees with sorting and search
// @Tags employee
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field" default(created_date)
// @Param sort_order query string false "Sort order (asc/desc)" default(asc)
// @Param search query string false "Search term"
// @Success 200 {object} EmployeeListResponse
// @Failure 500 {string} string "Server error"
// @Router /employees [get]
func GetEmployeeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_date"
	}
	sortOrder := query.Get("sort_order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	search := query.Get("search")

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = fmt.Sprintf(" WHERE (first_name_en ILIKE $%d OR last_name_en ILIKE $%d OR company_email ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	var totalRecords int
	DB.QueryRow("SELECT COUNT(*) FROM m_employee"+whereClause, args...).Scan(&totalRecords)

	offset := (page - 1) * pageSize
	mainQuery := fmt.Sprintf(`SELECT employee_id, employment_type, title, first_name_en, last_name_en, first_name_th, 
		last_name_th, nick_name_en, nick_name_th, phone_number, company_email, nationality, gender, tax_id, 
		birth_date, start_work_date, status, remark, department, position, photo, custom_attributes, created_by, 
		created_date, updated_by, updated_date, is_active FROM m_employee%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		whereClause, sortBy, strings.ToUpper(sortOrder), argIndex, argIndex+1)

	args = append(args, pageSize, offset)
	rows, err := DB.Query(mainQuery, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		var birthDate, startWorkDate, createdDate, updatedDate sql.NullTime
		var updatedBy sql.NullString

		rows.Scan(&emp.EmployeeID, &emp.EmploymentType, &emp.Title, &emp.FirstNameEN, &emp.LastNameEN,
			&emp.FirstNameTH, &emp.LastNameTH, &emp.NickNameEN, &emp.NickNameTH, &emp.PhoneNumber,
			&emp.CompanyEmail, &emp.Nationality, &emp.Gender, &emp.TaxID, &birthDate, &startWorkDate,
			&emp.Status, &emp.Remark, &emp.Department, &emp.Position, &emp.Photo, &emp.CustomAttributes,
			&emp.CreatedBy, &createdDate, &updatedBy, &updatedDate, &emp.IsActive)

		if birthDate.Valid {
			emp.BirthDate = birthDate.Time.Format("2006-01-02 15:04:05")
		}
		if startWorkDate.Valid {
			emp.StartWorkDate = startWorkDate.Time.Format("2006-01-02 15:04:05")
		}
		if createdDate.Valid {
			emp.CreatedDate = createdDate.Time.Format("2006-01-02 15:04:05")
		}
		if updatedBy.Valid {
			emp.UpdatedBy = updatedBy.String
		}
		if updatedDate.Valid {
			emp.UpdatedDate = updatedDate.Time.Format("2006-01-02 15:04:05")
		}

		employees = append(employees, emp)
	}

	totalPages := (totalRecords + pageSize - 1) / pageSize
	response := EmployeeListResponse{
		Data:       employees,
		Total:      totalRecords,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateEmployee godoc
// @Summary Update employee
// @Description Update employee information by ID
// @Tags employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body Employee true "Employee object"
// @Success 200 {object} Employee
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Server error"
// @Router /employee/{id} [put]
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employeeID := r.URL.Path[len("/api/employee/"):]
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE m_employee SET employment_type=$1, title=$2, first_name_en=$3, last_name_en=$4,
		first_name_th=$5, last_name_th=$6, nick_name_en=$7, nick_name_th=$8, phone_number=$9, company_email=$10,
		nationality=$11, gender=$12, tax_id=$13, birth_date=$14, start_work_date=$15, status=$16, remark=$17,
		department=$18, position=$19, photo=$20, custom_attributes=$21, updated_by=$22, updated_date=CURRENT_TIMESTAMP, 
		is_active=$23 WHERE employee_id=$24 RETURNING employee_id, employment_type, title, first_name_en, last_name_en, 
		first_name_th, last_name_th, nick_name_en, nick_name_th, phone_number, company_email, nationality, gender, tax_id, 
		birth_date, start_work_date, status, remark, department, position, photo, custom_attributes, created_by, created_date, 
		updated_by, updated_date, is_active`

	var updatedEmp Employee
	var birthDate, startWorkDate, createdDate, updatedDate sql.NullTime
	var updatedBy sql.NullString

	err := DB.QueryRow(query, employee.EmploymentType, employee.Title, employee.FirstNameEN, employee.LastNameEN,
		employee.FirstNameTH, employee.LastNameTH, employee.NickNameEN, employee.NickNameTH, employee.PhoneNumber,
		employee.CompanyEmail, employee.Nationality, employee.Gender, employee.TaxID, employee.BirthDate,
		employee.StartWorkDate, employee.Status, employee.Remark, employee.Department, employee.Position,
		employee.Photo, employee.CustomAttributes, employee.UpdatedBy, employee.IsActive, employeeID,
	).Scan(&updatedEmp.EmployeeID, &updatedEmp.EmploymentType, &updatedEmp.Title, &updatedEmp.FirstNameEN,
		&updatedEmp.LastNameEN, &updatedEmp.FirstNameTH, &updatedEmp.LastNameTH, &updatedEmp.NickNameEN,
		&updatedEmp.NickNameTH, &updatedEmp.PhoneNumber, &updatedEmp.CompanyEmail, &updatedEmp.Nationality,
		&updatedEmp.Gender, &updatedEmp.TaxID, &birthDate, &startWorkDate, &updatedEmp.Status,
		&updatedEmp.Remark, &updatedEmp.Department, &updatedEmp.Position, &updatedEmp.Photo,
		&updatedEmp.CustomAttributes, &updatedEmp.CreatedBy, &createdDate, &updatedBy, &updatedDate, &updatedEmp.IsActive)

	if err == sql.ErrNoRows {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if birthDate.Valid {
		updatedEmp.BirthDate = birthDate.Time.Format("2006-01-02 15:04:05")
	}
	if startWorkDate.Valid {
		updatedEmp.StartWorkDate = startWorkDate.Time.Format("2006-01-02 15:04:05")
	}
	if createdDate.Valid {
		updatedEmp.CreatedDate = createdDate.Time.Format("2006-01-02 15:04:05")
	}
	if updatedBy.Valid {
		updatedEmp.UpdatedBy = updatedBy.String
	}
	if updatedDate.Valid {
		updatedEmp.UpdatedDate = updatedDate.Time.Format("2006-01-02 15:04:05")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmp)
}

// DeleteEmployee godoc
// @Summary Delete employee
// @Description Delete employee by ID
// @Tags employee
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Server error"
// @Router /employee/{id} [delete]
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employeeID := r.URL.Path[len("/api/employee/"):]
	result, err := DB.Exec("DELETE FROM m_employee WHERE employee_id = $1", employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully", "id": employeeID})
}
