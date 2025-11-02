# Stage 1: Build
FROM golang:1.24.0-alpine AS builder

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
# Install Doppler CLI
RUN apk add --no-cache curl tzdata ca-certificates && \
    curl -Ls https://cli.doppler.com/install.sh | sh


ENV TZ=Europe/Madrid
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Copiem amb el nom correcte
COPY --from=builder /app/bin/actions-service /app/actions-service