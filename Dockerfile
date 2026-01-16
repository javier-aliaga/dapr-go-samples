# =========================
# Builder stage
# =========================
FROM golang:1.25-alpine AS builder

# Enable Go modules and turn on CGO if you need it (here we disable for static binary)
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install build tools (git often needed for go modules)
RUN apk add --no-cache git ca-certificates

# Cache dependencies first
COPY go.mod go.sum* ./
COPY vendor ./vendor

# If vendor is present, don't download modules; build using vendored deps
ENV GOFLAGS="-mod=vendor"

# Copy the rest of the source
COPY . .

# Build the binary
RUN go build -ldflags="-s -w" -o server ./main.go

# =========================
# Runtime stage
# =========================
FROM alpine:3.20

# Add CA certs for HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Expose HTTP port used by main.go
EXPOSE 8080

# Run the binary
CMD ["/app/server"]