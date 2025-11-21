# -----------------------------
# Stage 1: Build Go binary
# -----------------------------
FROM golang:1.24-alpine AS builder

LABEL maintainer="braiyenmassora@gmail.com"

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /usr/src/app

# Copy Go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary
RUN go build -o app .

# -----------------------------
# Stage 2: Minimal runtime image
# -----------------------------
FROM alpine:3.21.3

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl \
    && adduser -D user

# Set timezone to Jakarta
ENV TZ=Asia/Jakarta

# Set working directory
WORKDIR /usr/src/app

# Copy the binary from builder
COPY --from=builder --chown=user:user /usr/src/app/app ./app

# Switch to non-root user
USER user

# Expose the application port
EXPOSE 8080

# Optional: Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD curl -f http://localhost:8080/ || exit 1

# Run the application
CMD ["./app"]