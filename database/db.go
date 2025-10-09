package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
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

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error verifying connection to database:", err)
	}

	// Drop existing table if it exists (to fix column name case issues)
	_, _ = DB.Exec("DROP TABLE IF EXISTS m_employee")

	// Create employees table with lowercase column names
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS m_employee (
		employee_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		employment_type INT NOT NULL,
		title INT NOT NULL,
		first_name_en VARCHAR(50) NOT NULL,
		last_name_en VARCHAR(50) NOT NULL,
		first_name_th VARCHAR(50) NOT NULL,
		last_name_th VARCHAR(50) NOT NULL,
		nick_name_en VARCHAR(50) NOT NULL,
		nick_name_th VARCHAR(50) NOT NULL,
		phone_number VARCHAR(20) NOT NULL,
		company_email VARCHAR(320) NOT NULL,
		personal_email VARCHAR(320) NOT NULL,
		nationality VARCHAR(50) NOT NULL,
		gender INT NOT NULL,
		tax_id VARCHAR(13) NOT NULL,
		birth_date TIMESTAMP NOT NULL,
		start_work_date TIMESTAMP NOT NULL,
		status INT NOT NULL,
		remark TEXT NOT NULL,
		department VARCHAR(50) NOT NULL,
		position VARCHAR(50) NOT NULL,
		photo VARCHAR(256) NOT NULL,
		custom_attributes TEXT NOT NULL,
		created_by UUID NOT NULL,
		created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by UUID,
		updated_date TIMESTAMP,
		is_active BOOLEAN NOT NULL DEFAULT TRUE
	)`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("Database connection established and table created successfully")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
