# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files (we'll create these next)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flux-training-app .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/flux-training-app .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./flux-training-app"]
