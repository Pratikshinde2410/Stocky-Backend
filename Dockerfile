# syntax=docker/dockerfile:1

# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Ensure go.mod/go.sum include all needed deps
RUN go mod tidy

# Build the binary statically for a small final image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stocky ./cmd/server

# ---------- Runtime stage ----------
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/stocky /app/stocky

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/stocky"]


