# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o ducla-agent ./cmd/agent

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 ducla && \
    adduser -D -u 1000 -G ducla ducla

# Create necessary directories
RUN mkdir -p /opt/ducla/data /tmp/ducla /var/log/ducla /etc/ducla/tls && \
    chown -R ducla:ducla /opt/ducla /tmp/ducla /var/log/ducla /etc/ducla

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/ducla-agent .
COPY --from=builder /build/configs/agent.yaml /etc/ducla/agent.yaml

# Set ownership
RUN chown -R ducla:ducla /app

# Switch to non-root user
USER ducla

# Expose ports
EXPOSE 8080 8081 8443 9090

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the application
ENTRYPOINT ["/app/ducla-agent"]
CMD ["--config", "/etc/ducla/agent.yaml"]
