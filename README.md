# Lalan BE

Backend API for the Lalan Web Application – Admin Panel for Outdoor Hosters

## Requirements

- **Go** 1.24.4 or higher
- **PostgreSQL**

## Getting Started

**Environment Variables**

- `DATABASE_URL`: PostgreSQL connection string (required)

```bash
# Clone the repository
git clone https://github.com/ech0py/lalan-be.git
cd lalan-be

# Install dependencies
go mod download

# Configure environment
cp .env.dev .env

# Edit .env with your database credentials
DATABASE_URL='postgres://<USERNAME>:<PASSWORD>@<HOST>:<PORT>/<DB_NAME>?sslmode=require&search_path={schema_name}'
```

## Database Setup

Run the database migrations to create the required tables:

```bash
# Create hoster table
psql -U <username> -d <database> -f migrations/create_hoster.sql
# Create category table
psql -U <username> -d <database> -f migrations/create_categories.sql
# Create item table
psql -U <username> -d <database> -f migrations/create_item.sql
```

### Install Air (Hot Reload Tool)

Air provides automatic hot reload during development.

**Linux/Mac:**

```bash
# Install air
go install github.com/air-verse/air@latest

# Add Go bin to PATH (if not already added)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Verify installation
air -v
```

## Run Locally

### Development Mode (with hot reload)

```bash
make dev
```

### Alternative Methods

```bash
# Run without hot reload
go run ./cmd/main.go

# Build and run
go build -o main ./cmd/main.go
./main
```

## Project Structure

```
lalan-be/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/           # Database configuration
│   ├── handler/          # HTTP request handlers
│   ├── model/            # Data models
│   ├── repository/       # Database access layer
│   ├── response/         # HTTP response utilities
│   ├── route/            # Route definitions
│   └── service/          # Business logic layer
├── migrations/           # Database schema migrations
├── pkg/                  # Shared packages
├── .env.dev              # Development environment variables
├── go.mod                # Go module definition
└── go.sum                # Go module checksums
```

## API Documentation

### Authentication

- POST /v1/auth/register - Register a new hoster account
- POST /v1/auth/login - Authenticate hoster login

### Categories

- POST /v1/category/add - Create a new category

### Items

- POST /v1/item/add - Create a new item
