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

# Build executable amb el nom correcte
RUN go build -o /app/bin/actions-service ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Instal·lem dependencies i Doppler
RUN apk add --no-cache ca-certificates curl bash && \
    curl -Ls https://cli.doppler.com/install.sh | sh && \
    rm -rf /var/cache/apk/*

# Copiem amb el nom correcte
COPY --from=builder /app/bin/actions-service /app/actions-service

# No especifiquem CMD perquè ho fa el compose