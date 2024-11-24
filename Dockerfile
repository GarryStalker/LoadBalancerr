FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./
# GOOS - допустим мы билдим под линукс
# CGO_ENABLED - билдим без динамических либ, все нужные зависимости внутри бинаря
RUN CGO_ENABLED=0 GOOS=linux go build -o ./load-balancer ./cmd/server/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/load-balancer .

COPY config/local.yaml ./config/local.yaml

CMD ["./load-balancer", "--config=./config/local.yaml"]
