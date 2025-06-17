# ===========================
# Build Stage
# ===========================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# ===========================
# Production Stage
# ===========================
FROM alpine:3.18

WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary and necessary folders
COPY --from=builder /app/app ./
COPY --from=builder /app/view ./view
COPY --from=builder /app/public ./public

# Ensure upload directory exists and is writable
RUN mkdir -p /app/public/uploads \
    && chown -R appuser:appgroup /app/public \
    && chmod -R 755 /app/public

# Set non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./app"]
