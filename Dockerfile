# Stage 1: Build
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/actions-service ./cmd/api

# Stage 2: Runtime amb Doppler preinstal·lat
FROM dopplerhq/cli:3-alpine

WORKDIR /app

# Copiem el binari
COPY --from=builder /app/bin/actions-service /app/actions-service

