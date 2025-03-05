FROM golang:1.23 AS builder

WORKDIR /cmd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth_service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth_service .

EXPOSE 50051

CMD ["./auth_service"]