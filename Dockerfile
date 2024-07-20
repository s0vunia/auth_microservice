ARG CONFIG_PATH

FROM golang:1.20.3-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/crud_server .
COPY local.env .
COPY prod.env .

CMD ["./crud_server"]
