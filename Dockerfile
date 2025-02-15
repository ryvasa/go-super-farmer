# Base stage
FROM golang:1.23.4-alpine AS base

RUN apk add --no-cache git
WORKDIR /app

# Copy dependencies
COPY go.* ./
RUN go mod download

# Copy source code
COPY . .

# Development stage
FROM base AS development

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Use air for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM base AS production

RUN go build -o api cmd/main.go

# Runner stage
FROM alpine:latest AS runner

WORKDIR /app

COPY --from=production /app/api .

RUN chmod +x ./api

CMD ["./api"]
