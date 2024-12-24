
FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git

RUN mkdir /app

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN mkdir /build

RUN go build -o build/service_report cmd/report/main.go

# RUN chmod +x build/service_report

# CMD ["build/service_report"]

# Stage 2 running
FROM alpine:latest AS runner

RUN mkdir /app

COPY .env .

COPY --from=builder app/build/service_report /app

RUN chmod +x app/service_report

RUN mkdir /public

RUN mkdir /public/reports

CMD ["app/service_report"]
