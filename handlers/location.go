package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

// Geography represents the geography master data
type Geography struct {
	GeographyID int    `json:"geography_id"`
	Name        string `json:"name"`
}

// Province represents the province master data
type Province struct {
	ProvinceID     int    `json:"province_id"`
	ProvinceNameTH string `json:"province_name_th"`
	ProvinceNameEN string `json:"province_name_en"`
	GeographyID    int    `json:"geography_id"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	DeletedAt      string `json:"deleted_at,omitempty"`
}

// District represents the district master data
type District struct {
	DistrictID int    `json:"district_id"`
	NameTH     string `json:"name_th"`
	NameEN     string `json:"name_en"`
	ProvinceID int    `json:"province_id"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	DeletedAt  string `json:"deleted_at,omitempty"`
}

// SubDistrict represents the sub-district master data
type SubDistrict struct {
	SubDistrictID int    `json:"sub_district_id"`
	ZipCode       int    `json:"zip_code"`
	NameTH        string `json:"name_th"`
	NameEN        string `json:"name_en"`
	DistrictID    int    `json:"district_id"`
	Lat           string `json:"lat,omitempty"`
	Long          string `json:"long,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	DeletedAt     string `json:"deleted_at,omitempty"`
}

// GetGeographies godoc
// @Summary Get all geographies
// @Description Get list of all geographies
// @Tags location
// @Produce json
// @Success 200 {array} Geography
// @Failure 500 {string} string "Server error"
// @Router /geographies [get]
func GetGeographies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT geography_id, name 
		FROM m_geography 
		ORDER BY geography_id
	`

	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, "Error querying geographies: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var geographies []Geography
	for rows.Next() {
		var geography Geography
		if err := rows.Scan(&geography.GeographyID, &geography.Name); err != nil {
			http.Error(w, "Error scanning geography: "+err.Error(), http.StatusInternalServerError)
			return
		}
		geographies = append(geographies, geography)
	}

	if geographies == nil {
		geographies = []Geography{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geographies)
}

// GetProvinces godoc
// @Summary Get provinces
// @Description Get list of provinces, optionally filtered by geography_id
// @Tags location
// @Produce json
// @Param geography_id query int false "Geography ID to filter provinces"
// @Success 200 {array} Province
// @Failure 500 {string} string "Server error"
// @Router /provinces [get]
func GetProvinces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	geographyIDParam := r.URL.Query().Get("geography_id")

	var query string
	var rows *sql.Rows
	var err error

	if geographyIDParam != "" {
		var geographyID int
		geographyID, err = strconv.Atoi(geographyIDParam)
		if err != nil {
			http.Error(w, "Invalid geography_id parameter", http.StatusBadRequest)
			return
		}

		query = `
			SELECT province_id, province_name_th, province_name_en, geography_id, 
				   created_at, updated_at, deleted_at 
			FROM m_province 
			WHERE geography_id = $1 AND deleted_at IS NULL
			ORDER BY province_name_en
		`
		rows, err = DB.Query(query, geographyID)
	} else {
		query = `
			SELECT province_id, province_name_th, province_name_en, geography_id, 
				   created_at, updated_at, deleted_at 
			FROM m_province 
			WHERE deleted_at IS NULL
			ORDER BY province_name_en
		`
		rows, err = DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Error querying provinces: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var provinces []Province
	for rows.Next() {
		var province Province
		var createdAt, updatedAt, deletedAt sql.NullString

		if err := rows.Scan(
			&province.ProvinceID,
			&province.ProvinceNameTH,
			&province.ProvinceNameEN,
			&province.GeographyID,
			&createdAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			http.Error(w, "Error scanning province: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if createdAt.Valid {
			province.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			province.UpdatedAt = updatedAt.String
		}
		if deletedAt.Valid {
			province.DeletedAt = deletedAt.String
		}

		provinces = append(provinces, province)
	}

	if provinces == nil {
		provinces = []Province{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provinces)
}

// GetDistricts godoc
// @Summary Get districts
// @Description Get list of districts, optionally filtered by province_id
// @Tags location
// @Produce json
// @Param province_id query int false "Province ID to filter districts"
// @Success 200 {array} District
// @Failure 500 {string} string "Server error"
// @Router /districts [get]
func GetDistricts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	provinceIDParam := r.URL.Query().Get("province_id")

	var query string
	var rows *sql.Rows
	var err error

	if provinceIDParam != "" {
		var provinceID int
		provinceID, err = strconv.Atoi(provinceIDParam)
		if err != nil {
			http.Error(w, "Invalid province_id parameter", http.StatusBadRequest)
			return
		}

		query = `
			SELECT district_id, name_th, name_en, province_id, 
				   created_at, updated_at, deleted_at 
			FROM m_district 
			WHERE province_id = $1 AND deleted_at IS NULL
			ORDER BY name_en
		`
		rows, err = DB.Query(query, provinceID)
	} else {
		query = `
			SELECT district_id, name_th, name_en, province_id, 
				   created_at, updated_at, deleted_at 
			FROM m_district 
			WHERE deleted_at IS NULL
			ORDER BY name_en
		`
		rows, err = DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Error querying districts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var districts []District
	for rows.Next() {
		var district District
		var createdAt, updatedAt, deletedAt sql.NullString

		if err := rows.Scan(
			&district.DistrictID,
			&district.NameTH,
			&district.NameEN,
			&district.ProvinceID,
			&createdAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			http.Error(w, "Error scanning district: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if createdAt.Valid {
			district.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			district.UpdatedAt = updatedAt.String
		}
		if deletedAt.Valid {
			district.DeletedAt = deletedAt.String
		}

		districts = append(districts, district)
	}

	if districts == nil {
		districts = []District{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(districts)
}

// GetSubDistricts godoc
// @Summary Get sub-districts
// @Description Get list of sub-districts, optionally filtered by district_id
// @Tags location
// @Produce json
// @Param district_id query int false "District ID to filter sub-districts"
// @Success 200 {array} SubDistrict
// @Failure 500 {string} string "Server error"
// @Router /subdistricts [get]
func GetSubDistricts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	districtIDParam := r.URL.Query().Get("district_id")

	var query string
	var rows *sql.Rows
	var err error

	if districtIDParam != "" {
		var districtID int
		districtID, err = strconv.Atoi(districtIDParam)
		if err != nil {
			http.Error(w, "Invalid district_id parameter", http.StatusBadRequest)
			return
		}

		query = `
			SELECT sub_district_id, zip_code, name_th, name_en, district_id, 
				   lat, long, created_at, updated_at, deleted_at 
			FROM m_sub_district 
			WHERE district_id = $1 AND deleted_at IS NULL
			ORDER BY name_en
		`
		rows, err = DB.Query(query, districtID)
	} else {
		query = `
			SELECT sub_district_id, zip_code, name_th, name_en, district_id, 
				   lat, long, created_at, updated_at, deleted_at 
			FROM m_sub_district 
			WHERE deleted_at IS NULL
			ORDER BY name_en
		`
		rows, err = DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Error querying sub-districts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subDistricts []SubDistrict
	for rows.Next() {
		var subDistrict SubDistrict
		var lat, long, createdAt, updatedAt, deletedAt sql.NullString

		if err := rows.Scan(
			&subDistrict.SubDistrictID,
			&subDistrict.ZipCode,
			&subDistrict.NameTH,
			&subDistrict.NameEN,
			&subDistrict.DistrictID,
			&lat,
			&long,
			&createdAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			http.Error(w, "Error scanning sub-district: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if lat.Valid {
			subDistrict.Lat = lat.String
		}
		if long.Valid {
			subDistrict.Long = long.String
		}
		if createdAt.Valid {
			subDistrict.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			subDistrict.UpdatedAt = updatedAt.String
		}
		if deletedAt.Valid {
			subDistrict.DeletedAt = deletedAt.String
		}

		subDistricts = append(subDistricts, subDistrict)
	}

	if subDistricts == nil {
		subDistricts = []SubDistrict{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subDistricts)
}
