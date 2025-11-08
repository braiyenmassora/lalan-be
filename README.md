# Lalan BE

A Go-based backend API that provides administrative functionality for outdoor rental hosts, enabling management of hosts, categories, items, and bookings.

## Requirements

- Go 1.24.4 or higher
- PostgreSQL

## Getting Started

### Clone and Setup

```bash
git clone https://github.com/ech0py/lalan-be.git
cd lalan-be
go mod download
cp .env.dev .env
```

### Environment Configuration

Edit `.env` with database credentials:

```bash
DATABASE_URL='postgres://<USERNAME>:<PASSWORD>@<HOST>:<PORT>/<DB_NAME>?sslmode=require&search_path={schema_name}'
```

## Project Structure

```
lalan-be/
├── cmd/                    # Application entry point
├── internal/              # Core components
│   ├── config/            # Database configuration
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # JWT authentication
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

## Development Setup

### Install Air (Hot Reload Tool)

```bash
go install github.com/air-verse/air@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Run Locally

```bash
# Development mode with hot reload
make dev

# Alternative methods
go run ./cmd/main.go
# or
go build -o main ./cmd/main.go
./main
```

## API Documentation

### Authentication Flow

```http
POST /v1/auth/register  # Register new account
POST /v1/auth/login     # Login with existing account
```

### API Endpoints

#### Authentication

| Endpoint | Method | Auth |
|----------|--------|------|
| `/v1/auth/register` | POST | Public |
| `/v1/auth/login` | POST | Public |

#### Category

| Endpoint | Method | Auth |
|----------|--------|------|
| `/v1/category/list` | GET | Public |
| `/v1/category/detail?id={id}` | GET | Public |
| `/v1/category/add` | POST | Protected |
| `/v1/category/update?id={id}` | PUT | Protected |
| `/v1/category/delete?id={id}` | DELETE | Protected |

#### Item

| Endpoint | Method | Auth |
|----------|--------|------|
| `/v1/item/list` | GET | Public |
| `/v1/item/detail?id={id}` | GET | Public |
| `/v1/item/my-items` | GET | Protected |
| `/v1/item/add` | POST | Protected |
| `/v1/item/update?id={id}` | PUT | Protected |
| `/v1/item/delete?id={id}` | DELETE | Protected |

## Adding New Features

| Component | Description | Location |
|-----------|-------------|----------|
| Migration | Create or modify SQL files for database schema | `migrations/` |
| Model | Define data structures | `internal/model/` |
| Repository | Implement data access logic | `internal/repository/` |
| Service | Add business logic | `internal/service/` |
| Handler | Create HTTP handlers | `internal/handler/` |
| Routes | Register routes | `internal/route/` |
| Main | Initialize and register all components | `cmd/main.go` |
