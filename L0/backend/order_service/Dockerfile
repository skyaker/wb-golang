FROM golang:1.24.3-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/order_service ./cmd/order_service

FROM alpine:latest

COPY --from=builder /usr/local/bin/order_service /usr/local/bin/order_service

CMD ["order_service"]
