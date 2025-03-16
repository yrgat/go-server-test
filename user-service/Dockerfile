# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the service source code
COPY user-service/*.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/user-service .

# Expose port 3001
EXPOSE 3001

# Run the application
CMD ["./user-service"] 