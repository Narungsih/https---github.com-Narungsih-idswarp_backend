package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

// Department represents the department master data
type Department struct {
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	CreatedDate    string `json:"created_date,omitempty"`
	UpdatedDate    string `json:"updated_date,omitempty"`
}

// GetDepartments godoc
// @Summary Get all departments
// @Description Get list of all departments
// @Tags department
// @Produce json
// @Success 200 {array} Department
// @Failure 500 {string} string "Server error"
// @Router /departments [get]
func GetDepartments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT department_id, department_name, created_date, updated_date
		FROM r_department
		ORDER BY department_id
	`

	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, "Error querying departments: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var department Department
		if err := rows.Scan(
			&department.DepartmentID,
			&department.DepartmentName,
			&department.CreatedDate,
			&department.UpdatedDate,
		); err != nil {
			http.Error(w, "Error scanning department: "+err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, department)
	}

	if departments == nil {
		departments = []Department{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

// Position represents the position master data
type Position struct {
	PositionID   int    `json:"position_id"`
	DepartmentID int    `json:"department_id"`
	PositionName string `json:"position_name"`
	Acronym      string `json:"acronym"`
	CreatedDate  string `json:"created_date,omitempty"`
	UpdatedDate  string `json:"updated_date,omitempty"`
}

// GetPositions godoc
// @Summary Get positions
// @Description Get list of positions, optionally filtered by department_id
// @Tags department
// @Produce json
// @Param department_id query int false "Department ID to filter positions"
// @Success 200 {array} Position
// @Failure 500 {string} string "Server error"
// @Router /positions [get]
func GetPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	departmentIDParam := r.URL.Query().Get("department_id")

	var query string
	var rows *sql.Rows
	var err error

	if departmentIDParam != "" {
		var departmentID int
		departmentID, err = strconv.Atoi(departmentIDParam)
		if err != nil {
			http.Error(w, "Invalid department_id parameter", http.StatusBadRequest)
			return
		}

		query = `
			SELECT position_id, department_id, position_name, acronym, created_date, updated_date
			FROM r_position
			WHERE department_id = $1
			ORDER BY position_name
		`
		rows, err = DB.Query(query, departmentID)
	} else {
		query = `
			SELECT position_id, department_id, position_name, acronym, created_date, updated_date
			FROM r_position
			ORDER BY position_name
		`
		rows, err = DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Error querying positions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var positions []Position
	for rows.Next() {
		var position Position
		var createdDate, updatedDate sql.NullString

		if err := rows.Scan(
			&position.PositionID,
			&position.DepartmentID,
			&position.PositionName,
			&position.Acronym,
			&createdDate,
			&updatedDate,
		); err != nil {
			http.Error(w, "Error scanning position: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if createdDate.Valid {
			position.CreatedDate = createdDate.String
		}
		if updatedDate.Valid {
			position.UpdatedDate = updatedDate.String
		}

		positions = append(positions, position)
	}

	if positions == nil {
		positions = []Position{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(positions)
}
