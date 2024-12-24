
FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git

RUN mkdir /app

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN mkdir /build

RUN go build -o build/service_mail cmd/mail/main.go

# RUN chmod +x build/service_mail

# CMD ["build/service_mail"]

# Stage 2 running
FROM alpine:latest AS runner

RUN mkdir /app

COPY .env .

COPY --from=builder app/build/service_mail /app

RUN chmod +x app/service_mail

CMD ["app/service_mail"]
