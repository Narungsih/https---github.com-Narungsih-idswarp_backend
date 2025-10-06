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

	// Create employees table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS m_employee (
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
