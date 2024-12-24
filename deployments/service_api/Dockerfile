
FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git

RUN mkdir /app

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN mkdir /build

RUN go build -o build/service_api cmd/api/main.go

# RUN chmod +x build/service_api

# CMD ["build/service_api"]

# Stage 2 running
FROM alpine:latest AS runner

RUN mkdir /app

COPY .env .

RUN mkdir /pkg
RUN mkdir /pkg/auth
RUN mkdir /pkg/auth/casbin

COPY /pkg/auth/casbin/model.conf ./pkg/auth/casbin/

COPY /pkg/auth/casbin/policy.csv ./pkg/auth/casbin/

COPY --from=builder app/build/service_api /app

RUN chmod +x app/service_api

CMD ["app/service_api"]
