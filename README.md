# Lalan BE

A Go-based backend API for managing outdoor rental operations through an admin dashboard.  
Designed for scalability, maintainability, and clean architecture.

## Requirements

- Go 1.24.4 or higher  
- PostgreSQL  

## Getting Started

```bash
git clone https://github.com/ech0py/lalan-be.git
cd lalan-be
go mod download
```

## Environment Configuration

Set up the `.env.dev` file before running the application:

```bash
# PostgreSQL database connection (development)
DATABASE_URL='postgres://<USERNAME>:<PASSWORD>@<HOST>:<PORT>/<DB_NAME>?sslmode=require&search_path={schema_name}'

# JWT Secret Key openssl rand -base64 32
JWT_SECRET=<your_secret_key>

# Application Port
APP_PORT=8080

```

## Project Structure

```
lalan-be/
├── cmd/                    # Application entry point
├── internal/               # Core logic and modules
│   ├── config/             # App and database configuration
│   ├── handler/            # HTTP request handlers
│   ├── middleware/         # Authentication and middleware logic
│   ├── model/              # Data models
│   ├── repository/         # Database operations
│   ├── response/           # Response formatting utilities
│   ├── route/              # Route definitions
│   └── service/            # Business logic layer
├── migrations/             # Database migrations
├── pkg/                    # Shared helper packages
├── .env.dev                # Environment configuration (development)
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

## Development Setup

### Install Air (Live Reload)

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

## API Documentation

### Authentication Flow

```http
POST /v1/auth/register   # Register a new account
POST /v1/auth/login      # Login with existing credentials
```

### API Endpoints

#### Authentication

| Endpoint              | Method | Access   |
|-----------------------|--------|----------|
| `/v1/auth/register`   | POST   | Public   |
| `/v1/auth/login`      | POST   | Public   |

#### Category

| Endpoint                      | Method | Access    |
|-------------------------------|--------|-----------|
| `/v1/category/list`           | GET    | Public    |
| `/v1/category/detail?id={id}` | GET    | Public    |
| `/v1/category/add`            | POST   | Protected |
| `/v1/category/update?id={id}` | PUT    | Protected |
| `/v1/category/delete?id={id}` | DELETE | Protected |

#### Item

| Endpoint                   | Method | Access    |
|-----------------------------|--------|-----------|
| `/v1/item/list`             | GET    | Public    |
| `/v1/item/detail?id={id}`   | GET    | Public    |
| `/v1/item/my-items`         | GET    | Protected |
| `/v1/item/add`              | POST   | Protected |
| `/v1/item/update?id={id}`   | PUT    | Protected |
| `/v1/item/delete?id={id}`   | DELETE | Protected |

---

## Adding New Features

| Component  | Description                              | Location               |
|-------------|------------------------------------------|------------------------|
| Migration   | Manage database schema changes           | `migrations/`          |
| Model       | Define data structures                   | `internal/model/`      |
| Repository  | Implement database access logic          | `internal/repository/` |
| Service     | Handle business logic                    | `internal/service/`    |
| Handler     | Create HTTP request handlers             | `internal/handler/`    |
| Routes      | Register and manage API routes           | `internal/route/`      |
| Main        | Initialize and link all core components  | `cmd/main.go`          |

---

## Code Commenting Guidelines

```bash
{
  "initial": "Review the entire codebase to understand its purpose before adding comments.",
  "after": "Remove outdated or unnecessary comments to maintain a clean and consistent codebase.",
  "new": "Organize code sections according to Golang best practices, then add concise technical comments for each function, method, constant, struct, or significant part of the code. Use block comment format /* ... */. Split long comments into multiple lines for readability. Focus comments on explaining purpose and expected outcomes, avoid restating obvious implementation details. Ensure consistency and clarity in all comments throughout the project."
}
```
