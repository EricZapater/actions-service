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
RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg && \
    curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list && \
    apt-get update && \
    apt-get -y install doppler

# Copiem amb el nom correcte
COPY --from=builder /app/bin/actions-service /app/actions-service

# No especifiquem CMD perquè ho fa el compose