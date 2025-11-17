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

## Code Commenting Guidelines

```
{
  "task": "refactor_go_existing_file_as_agent_strict_no_logic_change",
  "goal": "Refactor file Go agar struktur super rapi dan mudah dibaca dengan komentar bahasa Indonesia yang konsisten, ringkas, dan langsung bisa dipahami junior — TANPA MENGUBAH LOGIKA KODE SAMA SEKALI.",
  "requirements": {
    "analysis": [
      "Identifikasi semua deklarasi: package, import, const, var, type, struct, interface, method (receiver), func (termasuk init), dan main."
    ],
    "comment_cleanup": "Hapus 100% komentar lama (// maupun /* */). Semua akan diganti dengan komentar baru yang jauh lebih baik.",
    "strict_ordering": [
      "package declaration",
      "imports (grouped std/third-party/local & sorted alphabetically)",
      "const blocks",
      "var blocks",
      "init() functions",
      "type declarations + methods (dikelompokkan per receiver type)",
      "interfaces",
      "private functions (sorted alphabetically)",
      "public functions (sorted alphabetically)",
      "func main() → WAJIB di baris paling akhir"
    ],
    "prepend_comments": {
      "language": "Indonesian_only",
      "style": "block_comment_exactly",
      "position": "Tepat satu baris kosong di atas setiap elemen/grup yang wajib dikomentari",
      "mandatory_for": [
        "every const block",
        "every var block",
        "every type declaration",
        "every method group (per receiver)",
        "every interface",
        "every standalone function (termasuk init)"
      ],
      "template_strict": "/*\n[Tujuan elemen ini dalam 1 kalimat singkat].\n[Hasil atau konteks penggunaannya dalam 1 kalimat singkat].\n*/",
      "rules": [
        "Maksimal 2 kalimat, maksimal 3 baris total.",
        "Bahasa Indonesia baku tapi santai, langsung to the point.",
        "Dilarang keras pakai 'Ini adalah', 'Fungsi ini', 'Berfungsi untuk' — langsung inti saja."
      ],
      "examples": {
        "const": "/*\nKode status HTTP standar untuk seluruh aplikasi.\nDipakai di handler dan middleware respons.\n*/",
        "var": "/*\nKoneksi pool PostgreSQL yang dipakai global.\nDi-initialize saat startup dan tidak pernah di-close manual.\n*/",
        "type": "/*\nUserService menangani semua logika bisnis user.\nMenggunakan repository dan email service.\n*/",
        "method": "/*\nLogin memverifikasi credential dan mengeluarkan JWT.\nMengembalikan error jika password salah atau akun terkunci.\n*/",
        "interface": "/*\nRepository mendefinisikan kontrak akses data.\nWajib diimplementasikan oleh semua storage (SQL, Mongo, dll).\n*/",
        "function": "/*\nvalidateRequest memeriksa header, body, dan query param.\nMengembalikan error 400 jika ada field wajib yang kosong.\n*/"
      }
    },
    "sql_formatting": "Ubah semua query SQL menjadi raw string literal dengan backtick (`...`) dan indentasi 4 spasi. Placeholder parameter TIDAK BOLEH berubah.",
    "allowed_changes": [
      "Reorder deklarasi sesuai strict_ordering",
      "Grouping method berdasarkan receiver type",
      "Tambah/hapus/ganti komentar",
      "Formatting, spacing, line break, alignment",
      "Ubah string SQL biasa → backtick multiline"
    ],
    "ABSOLUTELY_FORBIDDEN": [
      "Mengubah apapun yang mempengaruhi logika eksekusi program",
      "Mengganti nama variabel, fungsi, type, method, struct field (kecuali typo fatal yang jelas salah eja)",
      "Menambah, menghapus, atau memindah baris kode fungsional",
      "Mengubah urutan ekspresi, kondisi, loop, atau return value",
      "Mengubah signature function/method yang bersifat public (huruf kapital)",
      "Mengubah nilai default const/var (kecuali formatting)",
      "Mengganti operator, tipe data, atau pointer/reference"
    ],
    "agent_behavior": "Kamu adalah agen yang HANYA boleh melakukan cosmetic & documentation changes. Jika ragu-ragu apakah suatu perubahan akan mengubah logika, JANGAN LAKUKAN. Kerjakan bertahap jika file besar, tapi setiap diff yang kamu kirim HARUS 100% aman dan sesuai aturan di atas.",
    "final_rule": "Lebih baik kamu melewatkan satu komentar daripada mengubah satu baris pun yang bisa mempengaruhi runtime behavior. Safety first."
  }
}
```
