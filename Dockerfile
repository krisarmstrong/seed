# Build stage for Go backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev libpcap-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o netscope ./cmd/netscope

# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package*.json ./
RUN npm ci

# Copy frontend source
COPY web/ ./

# Build frontend
RUN npm run build

# Final runtime image
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    libpcap \
    ca-certificates \
    tzdata

# Copy binary from builder
COPY --from=backend-builder /app/netscope .

# Copy frontend dist
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy default config
COPY configs/ ./configs/

# Create non-root user (but allow cap_net_raw)
RUN adduser -D -u 1000 netscope

# Expose port
EXPOSE 8080 8443

# Default command (can be overridden)
CMD ["/app/netscope"]
