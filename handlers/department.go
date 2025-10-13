package handlers

import (
	"encoding/json"
	"net/http"
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
