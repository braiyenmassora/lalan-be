# -----------------------------------------------------------------------------
# Stage 1: Build static Go binary
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

LABEL maintainer="braiyenmassora@gmail.com"

# Install only required build dependencies
RUN apk add --no-cache ca-certificates gcc musl-dev

# Set working directory
WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build fully static binary (no CGO, smallest size, no libc dependencies)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build \
    -ldflags="-s -w" \
    -trimpath \
    -o app \
    ./cmd

# -----------------------------------------------------------------------------
# Stage 2: Minimal runtime image (distroless-like using scratch)
# -----------------------------------------------------------------------------
FROM scratch

# Copy CA certificates for HTTPS/TLS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data (Jakarta)
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
ENV TZ=Asia/Jakarta

# Copy only the compiled binary
COPY --from=builder /app/app /app

# Run as non-root by default (scratch uses UID 65534/nobody)
# No shell, no users, no package manager = maximum security

# Expose application port
EXPOSE 8080

# Health check (assuming you have /health endpoint in your app)
# If not, replace with: CMD ["/app", "version"] or remove this block
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD ["/app", "health"] || exit 1

# Run the binary
ENTRYPOINT ["/app"]