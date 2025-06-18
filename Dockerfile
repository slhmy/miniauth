# Multi-stage Dockerfile for MiniAuth
# Stage 1: Build frontend (React/Vite)
FROM node:22-alpine AS frontend-builder

WORKDIR /app/website

# Copy package files first for better caching
COPY website/package.json website/pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && \
    pnpm install --frozen-lockfile

# Copy source code
COPY website/ ./

# Build the frontend
RUN pnpm run build && \
    # Remove source maps and unnecessary files to reduce size
    find dist -name "*.map" -delete

# Stage 2: Build backend (Go)
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install build dependencies including gcc for CGO
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the Go application with optimizations
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -a -installsuffix cgo -o miniauth .

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite tzdata wget && \
    # Create non-root user for security
    addgroup -g 1001 miniauth && \
    adduser -D -s /bin/sh -u 1001 -G miniauth miniauth

WORKDIR /app

# Copy the Go binary from builder stage
COPY --from=backend-builder /app/miniauth .

# Copy frontend build files
COPY --from=frontend-builder /app/website/dist ./website/dist

# Copy any necessary configuration files
COPY --from=backend-builder /app/api ./api

# Create directory for database and set proper permissions
RUN mkdir -p /data && \
    chown -R miniauth:miniauth /app /data

# Switch to non-root user
USER miniauth

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080 \
    DB_TYPE=sqlite \
    DB_PATH=/data/miniauth.db \
    GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./miniauth"]
