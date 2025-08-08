# build
FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o user-server ./cmd/grpc_server/main.go

# runtime
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/user-server .

CMD ["./user-server"]