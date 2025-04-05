# Use Go 1.24 base image (Debian-based for compatibility)
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o main .

# Final minimal image
FROM debian:bullseye-slim

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy public assets if needed
# COPY --from=builder /app/public ./public

# Expose port
EXPOSE 8080

# Run the app
CMD ["./main"]
