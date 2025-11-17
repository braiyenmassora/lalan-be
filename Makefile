# Menyatakan target phony untuk menghindari konflik dengan file bernama sama
.PHONY: dev build run clean install-air

# Menjalankan server development dengan hot reload menggunakan air
dev:
	~/go/bin/air

# Membangun binary aplikasi ke direktori tmp
build:
	go build -o ./tmp/main ./cmd/main.go

# Menjalankan aplikasi langsung tanpa membangun
run:
	go run ./cmd/main.go

# Membersihkan file build sementara dan direktori
clean:
	rm -rf ./tmp

# Menginstall tool air untuk fitur hot reload
install-air:
	go install github.com/air-verse/air@latest