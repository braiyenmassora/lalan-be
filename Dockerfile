# -----------------------------
# Stage 1: Build static Go binary
# -----------------------------
FROM golang:1.24-alpine AS builder

LABEL maintainer="braiyenmassora@gmail.com"

RUN apk add --no-cache ca-certificates gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -trimpath -o app ./cmd

# -----------------------------
# Stage 2: Runtime image (distroless — kecil, aman, sudah ada timezone)
# -----------------------------
FROM gcr.io/distroless/static-debian12

# Copy binary
COPY --from=builder /app/app /app

# Timezone otomatis Jakarta (distroless sudah support TZ)
ENV TZ=Asia/Jakarta

# Non-root user (sudah default)
EXPOSE 8080

# Healthcheck (kalau punya /health endpoint)
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD ["/app", "health"] || exit 1

ENTRYPOINT ["/app"]