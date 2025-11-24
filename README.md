# Lalan BE

A Go-based backend API for managing outdoor rental operations through an admin dashboard.  
Designed for scalability, maintainability, and clean architecture.

## Requirements

- Go 1.24.4 or higher
- PostgreSQL

## Getting Started

```bash
git clone https://github.com/braiyenmassora/lalan-be.git
cd lalan-be
go mod download
```

## Environment Configuration

Set up the `.env.dev` file before running the application:

```bash
# JWT Secret Key
# Generate with: openssl rand -base64 32
JWT_SECRET=""

# Application Environment (dev or prod)
APP_ENV=dev

# Application Port
APP_PORT=8080

# PostgreSQL database connection (development)
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
```

## Project Structure

```
lalan-be/
├── cmd/                        # Application entry point
├── internal/                   # Core logic and modules
│   ├── config/                 # App and database configuration
│   ├── features/               # Feature-based modules
│   │   ├── admin/              # Admin-specific features
│   │   │   ├── handler.go      # Admin HTTP handlers
│   │   │   ├── repository.go   # Admin database operations
│   │   │   ├── route.go        # Admin route definitions
│   │   │   └── service.go      # Admin business logic
│   │   ├── hoster/             # Hoster-specific features
│   │   │   ├── handler.go      # Hoster HTTP handlers
│   │   │   ├── repository.go   # Hoster database operations
│   │   │   ├── route.go        # Hoster route definitions
│   │   │   └── service.go      # Hoster business logic
│   │   └── public/             # Public features (no auth required)
│   │       ├── handler.go      # Public HTTP handlers
│   │       ├── repository.go   # Public database operations
│   │       ├── route.go        # Public route definitions
│   │       └── service.go      # Public business logic
│   ├── middleware/             # Authentication and middleware logic
│   ├── model/                  # Data models
│   ├── repository/             # Shared repository interfaces
│   ├── response/               # Response formatting utilities
│   ├── route/                  # Shared route setup
│   └── service/                # Shared service interfaces
├── migrations/                 # Database migrations
├── pkg/                        # Shared helper packages
├── .env.dev                    # Environment configuration (development)
├── go.mod                      # Go module definition
└── go.sum                      # Go module checksums
```

## Run Locally

### Install Air

```bash
go install github.com/air-verse/air@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Run Application

```bash
# Run in development mode with live reload
make dev

# Or run manually
go run ./cmd/main.go

# Or build and execute binary
go build -o main ./cmd/main.go
./main
```

## Adding New Features

| Component  | Description                              | Location               |
|------------|------------------------------------------|------------------------|
| Migration  | Manage database schema changes           | `migrations/`          |
| Model      | Define data structures                   | `internal/model/`      |
| Repository | Implement database access logic          | `internal/repository/` |
| Service    | Handle business logic                    | `internal/service/`    |
| Handler    | Create HTTP request handlers             | `internal/handler/`    |
| Routes     | Register and manage API routes           | `internal/route/`      |
| Main       | Initialize and link all core components  | `cmd/main.go`          |
