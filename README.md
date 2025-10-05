# Employee Management API

A RESTful API for managing employees with PostgreSQL database, built with Go.

## Features

- ✅ Create new employees
- ✅ Get employee by ID
- ✅ PostgreSQL database integration
- ✅ Swagger UI documentation
- ✅ CORS enabled
- ✅ Environment variables configuration

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Git

## Setup

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd backend
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=IDS-warp
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
```

### 4. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080`
