# Start from the official Golang image
FROM golang:1.21-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port (change if your app uses another)
EXPOSE 8080

# Run the executable
CMD ["./main"]
