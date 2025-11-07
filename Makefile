# Declare phony targets to avoid conflicts with files of the same name
.PHONY: dev build run clean install-air

# Run development server with hot reload using air
dev:
	~/go/bin/air

# Build the application binary to tmp directory
build:
	go build -o ./tmp/main ./cmd/main.go

# Run the application directly without building
run:
	go run ./cmd/main.go

# Clean temporary build files and directories
clean:
	rm -rf ./tmp

# Install air tool for hot reload functionality
install-air:
	go install github.com/air-verse/air@latest
