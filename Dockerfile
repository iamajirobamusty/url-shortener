# ==========================================
# STAGE 1: Build the Static Binary
# ==========================================
FROM golang:1.26-alpine AS builder

# Install build dependencies (git, ca-certificates, tzdata)
RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy dependency manifests first to leverage Docker's layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the Go application as a static binary 
# CGO_ENABLED=0 disables C bindings for cross-compilation safety
# GOOS=linux forces compilation targeting runtime architectures
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o url-shortener cmd/server/main.go

# ==========================================
# STAGE 2: Create the Lightweight Runtime Image
# ==========================================
FROM alpine:3.19 AS final

# Secure runtime configuration
RUN apk update && apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 appuser

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/url-shortener .

# Copy the environment file template if required by deployment runtimes
# Note: Production values are injected via platform environment controls!
COPY .env .env

# Switch away from root to an unprivileged worker profile for enhanced container hardening
USER appuser

# Expose target routing channels
EXPOSE 8080

# Kick off engine execution
CMD ["./url-shortener"]