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
# Install Doppler CLI
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler

# Copiem amb el nom correcte
COPY --from=builder /app/bin/actions-service /app/actions-service

# No especifiquem CMD perquè ho fa el compose