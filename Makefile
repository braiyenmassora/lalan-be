.PHONY: dev build run clean install-air

dev:
	~/go/bin/air

build:
	go build -o ./tmp/main ./cmd/main.go

run:
	go run ./cmd/main.go

clean:
	rm -rf ./tmp

install-air:
	go install github.com/air-verse/air@latest
