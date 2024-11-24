FROM golang:1.22
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o ./load-balancer ./cmd/server/
CMD ["./load-balancer", "--config=./config/local.yaml"]