# Build stage
FROM golang:1.24.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for private modules)
RUN apk add --no-cache git

# Copy go mod files and download dependencies first (for better caching)
COPY go.mod ./
COPY go.su[m] ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001 -G appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appuser main

# Switch to non-root user
USER appuser

# Expose port 8000 (default port from main.go)
EXPOSE 8000

# Run the binary
CMD ["./main"]