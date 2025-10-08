# Stage 1: Build
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Instal·lem depèndencies necessàries
RUN apk add --no-cache git ca-certificates

# Copiem go.mod i go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiem la resta del codi
COPY . .

# Build executable
RUN go build -o /app/bin/server ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Certs per HTTPS
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/server /app/server

RUN curl -sLf https://cli.doppler.com/install.sh | sh

# Executem el binari
CMD ["/app/server"]
