# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
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

# Create non-root user to match Helm chart security context
RUN addgroup -g 10001 -S appuser && \
    adduser -u 10001 -S appuser -G appuser

# Create necessary directories and set permissions
RUN mkdir -p /app && \
    chown -R appuser:appuser /app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appuser /app/flux-training-app /app/flux-training-app

# Switch to non-root user
USER 10001:10001

# Expose port 8080
EXPOSE 8080

# Use ENTRYPOINT instead of CMD for better compatibility
ENTRYPOINT ["/app/flux-training-app"]