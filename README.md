# Lalan BE

Backend API for the Lalan Web Application – Admin Panel for Outdoor Hosters

## Requirements

- **Go** 1.24.4 or higher
- **PostgreSQL**

## Getting Started

### 1. Clone & Install Dependencies

```bash
# Clone the repository
git clone https://github.com/ech0py/lalan-be.git
cd lalan-be

# Install dependencies
go mod download
```

### 2. Configure Environment

Copy the example environment file and configure your database settings:

```bash
cp .env.local .env
# Edit .env with your database credentials
```

### 3. Setup Database

Run the database migration:

```bash
# Execute migration file
psql -U your_user -d your_database -f migrations/create_hoster.sql
```

### 4. Install Air (Hot Reload Tool)

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

## Running the Application

### Development Mode (with hot reload)

```bash
make dev
```

### Alternative Methods

```bash
# Run with air directly
air

# Run without hot reload
go run ./cmd/main.go

# Build and run
make build
./tmp/main
```

## Project Structure

```
lalan-be/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/           # Configuration
│   ├── handler/          # HTTP handlers
│   ├── helper/           # Helper functions
│   ├── model/            # Data models
│   ├── repository/       # Database layer
│   ├── route/            # Route definitions
│   └── service/          # Business logic
├── migrations/           # Database migrations
├── pkg/                  # Public packages
├── tmp/                  # Temporary build files
├── .env.dev              # Development environment variables
├── .air.toml             # Air configuration for hot reload
└── Makefile              # Build and run commands
```

## Available Commands

```bash
make dev          # Run development server with hot reload
make build        # Build the application
make run          # Run the application
make clean        # Clean temporary files
make install-air  # Install air tool
```

## License

MIT
