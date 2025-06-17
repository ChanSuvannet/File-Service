# First stage: build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (for go modules that might need it)
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build binary
RUN go build -o app .

# Final stage: production-ready image using distroless
FROM gcr.io/distroless/base

WORKDIR /app

# Copy built binary
COPY --from=builder /app/app .

# Copy necessary folders
COPY view /app/view
COPY public /app/public

# Expose the port
EXPOSE 8080

# Start the application
CMD ["/app/app"]
