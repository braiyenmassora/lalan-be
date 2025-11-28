# Lalan BE

API for managing outdoor rental operations used by admin, hoster, and customer.

## Tech Stack

| Field            | Value                                              |
|------------------|----------------------------------------------------|
| Language         | Go (1.24.4+)                                       |
| Database         | PostgreSQL (via Supabase)                          |
| Storage Bucket   | Supabase Bucket                                    |
| Payment Gateway  | Xendit                                             |

## Getting Started

Clone the project:

```bash
git clone https://github.com/braiyenmassora/lalan-be.git
```

Go to the project directory and install dependencies:

```bash
go mod download
```

Set up `.env.dev`:

```bash
# Application Environment
APP_ENV=
APP_PORT=

# Database Configuration (Supabase Pooler)
DB_HOST=
DB_NAME=
DB_PASSWORD=
DB_PORT=
DB_USER=
DB_SSL_MODE=

# JWT Configuration - IMPORTANT: Without quotes!
JWT_SECRET=

# Redis Configuration (currently commented out in code)
REDIS_URL=

# CORS Configuration (optional, defaults to *)
# ALLOWED_ORIGIN=

# Supabase Storage
STORAGE_ACCESS_KEY=
STORAGE_SECRET_KEY=
STORAGE_ENDPOINT=
STORAGE_REGION=
STORAGE_BUCKET=
STORAGE_PROJECT_ID=
STORAGE_DOMAIN=
```

Start the server:

```bash
make dev
# or
go run ./cmd/main.go
```

## Architecture

```
lalan-be/
├── cmd/                        # entrypoint (main.go)
├── internal/
│   ├── config/                 # configuration (DB, env, redis)
│   ├── domain/                 # entity / model DB (single source of truth)
│   ├── dto/                    # request & response DTO (API contracts)
│   ├── features/               # vertical slices (per actor / feature)
│   │   ├── admin/
│   │   ├── auth/
│   │   ├── customer/
│   │   ├── hoster/
│   │   └── public/
│   ├── middleware/             # auth, logging, CORS, etc
│   ├── response/               # standard API response
│   └── utils/                  # helper (upload, time, etc)
├── migrations/                 # SQL migrations
└── .env.dev
```

## Feature Guide

| Step | Description |
|------|-------------|
| 1    | Definisikan model di internal/domain jika membutuhkan tabel baru. |
| 2    | Buat DTO untuk request/response di internal/dto/ untuk menjaga kontrak API. |
| 3    | Buat folder fitur: `internal/features/[actor]/[feature]/` dengan minimal file: handler.go (HTTP layer), service.go (business logic), repository.go (data access), route.go (endpoint registration). |
| 4    | Tambahkan migrasi SQL di migrations/ bila perlu. |
| 5    | Sertakan tes dan dokumentasi endpoint saat menyelesaikan fitur. |
